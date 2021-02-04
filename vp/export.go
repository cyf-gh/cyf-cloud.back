package vp

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
)

func FillCellInfo2( f *excelize.File, d *Data, i, info2Index int ) {
	for _, j1 := range d.BasicInfo2[info2Index].Childs {
		v, _ := f.GetCellValue("Sheet1", "C"+strconv.Itoa( i ) )
		if v == "" { continue }
		for _, c := range j1.Childs {
			v = strings.Replace(v, " ", "", -1)
			// println( v, c.K )
			if c.K == "" {
				f.SetCellValue("Sheet1", "D29", c.V )
			}
			if v == c.K {
				f.SetCellValue("Sheet1", "D"+strconv.Itoa( i ), string( c.V ) )
				// println( c.V )
				break
			}
		}
	}
}

func FillCellInfo1( f *excelize.File, d *Data, i, info1Index int ) {
	for _, j1 := range d.BasicInfo[info1Index].Childs {
		// println(j1.V)
		v, _ := f.GetCellValue("Sheet1", "A"+strconv.Itoa( i ) )
		if v == j1.K {
			f.SetCellValue("Sheet1", "B"+strconv.Itoa( i ), j1.V)
			break
		}
	}
}

func Export( data, templatePath , savePath string, imgBytes []byte ) (e error) {
	f, e := excelize.OpenFile( templatePath )
	var d Data
	json.Unmarshal( []byte( data ), &d )
	// info 1
	for i := 2; i < 16; i++ {
		FillCellInfo1( f, &d, i, 0 )
	}
	for i := 17; i < 30; i++ {
		FillCellInfo1( f, &d, i, 1 )
	}
	// info 2
	for i := 2; i < 16; i++ {
		FillCellInfo2( f, &d, i, 0 )
	}
	for i := 17; i < 30; i++ {
		FillCellInfo2( f, &d, i, 1 )
	}
	// 填写流程
	curLine := 2
	for _, p := range d.Progress {
		// 填写主流程名字
		f.SetCellValue("Sheet1", "F"+strconv.Itoa( curLine ), p.Name )
		for _, c := range p.Childs {
			f.SetCellValue("Sheet1", "G"+strconv.Itoa( curLine ), c.Name )
			f.SetCellValue("Sheet1", "H"+strconv.Itoa( curLine ), c.Percent + "%" )
			f.SetCellValue("Sheet1", "I"+strconv.Itoa( curLine ), c.Date )
			curLine++
		}
	}
	// 插入图片
	e = f.AddPictureFromBytes("Sheet1", "A33", "","chart", ".png", imgBytes )
	// 根据指定路径保存文件
	e = f.SaveAs( savePath )
	return
}