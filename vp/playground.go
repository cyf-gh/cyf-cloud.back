// virtual progress chart
package vp

import (
	"../cc/err"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strconv"
)

func Test() {
	f, e := excelize.OpenFile("C:/.repos/cyf-cloud.back/.vp/template.xlsx"); err.Assert( e )
	var d Data
	json.Unmarshal( []byte(`{"chartProps":{"XDelta":20,"nodes":[{"name":"确定方案","x":72.71739130434783,"y":455},{"name":"收集资料","x":154.23913043478262,"y":432.5},{"name":"交支行","x":238.47826086956522,"y":387.5},{"name":"退回","x":243.91304347826087,"y":405.5},{"name":"交分行","x":485.76086956521743,"y":365},{"name":"评级完成","x":570,"y":342.5},{"name":"交支行","x":708.5869565217392,"y":297.5},{"name":"退回","x":735.7608695652175,"y":320},{"name":"交分行","x":738.4782608695652,"y":275},{"name":"上会","x":741.195652173913,"y":252.49999999999997},{"name":"授信批复","x":743.9130434782609,"y":230},{"name":"交支行","x":904.2391304347826,"y":185},{"name":"退回","x":988.4782608695651,"y":221},{"name":"提款批复","x":1070,"y":149},{"name":"退回","x":1235.7608695652173,"y":135.5},{"name":"放行","x":1320,"y":50}],"totalDays":460,"h":600,"w":1400,"TagAng":-1.0471975511965976},"marking":[{"Name":"营销","Childs":[{"Name":"拜访","Percent":"6"},{"Name":"确定方案","Percent":"10"},{"Name":"收集资料","Percent":"15"}]},{"Name":"评级","Childs":[{"Name":"撰写评级","Percent":"20"},{"Name":"交支行","Percent":"25"},{"Name":"退回","Percent":"21"},{"Name":"交分行","Percent":"30"},{"Name":"评级完成","Percent":"35"}]},{"Name":"授信","Childs":[{"Name":"撰写授信","Percent":"40"},{"Name":"交支行","Percent":"45"},{"Name":"退回","Percent":"40"},{"Name":"交分行","Percent":"50"},{"Name":"上会","Percent":"55"},{"Name":"授信批复","Percent":"60"}]},{"Name":"提款","Childs":[{"Name":"撰写提款","Percent":"61"},{"Name":"交支行","Percent":"70"},{"Name":"退回","Percent":"62"},{"Name":"提款批复","Percent":"78"}]},{"Name":"放款","Childs":[{"Name":"提交","Percent":"80"},{"Name":"退回","Percent":"81"},{"Name":"放行","Percent":"100"}]}],"basicInfo":[{"Title":"客户营销/需求","Childs":[{"k":"客户名","v":"1"},{"k":"种类","v":"2"},{"k":"金额(万元)","v":"3"},{"k":"特殊套用","v":"4"},{"k":"价格","v":"5"},{"k":"时间节点","v":"6"},{"k":"项目提要","v":"7"},{"k":"他行竞争/续作情况","v":"8"},{"k":"可能风险点","v":"9"},{"k":"材料清单","v":"10"}]},{"Title":"风险审批","Childs":[{"k":"评级撰写质量","v":"1"},{"k":"调查报告质量","v":"2"},{"k":"需与分行确认的问题","v":"3"},{"k":"营销部门未发现的瑕疵","v":"4"},{"k":"其他风险-非实质性","v":"5"},{"k":"其他风险-实质性","v":"6"},{"k":"补充材料要求","v":"7"}]}],"basicInfo2":[{"Title":"部门尽调","Childs":[{"Title":"尽调发现风险","Childs":[{"k":"非实质性","v":"无"},{"k":"实质性","v":"无"}]},{"Title":"需协调问题","Childs":[{"k":"柜面&账务","v":"无"},{"k":"授信政策","v":"无"},{"k":"其他","v":"无"}]},{"Title":"风险审批发现问题是否已解决","Childs":[{"k":"非实质性","v":"无"},{"k":"实质性","v":"无"}]}]},{"Title":"流程干预","Childs":[{"Title":"行长室意见","Childs":[{"k":"非实质性风险","v":"1"},{"k":"实质性风险","v":"2"},{"k":"其他","v":"3"}]},{"Title":"时长统计","Childs":[{"k":"营销评级","v":"123"},{"k":"营销授信","v":"23"},{"k":"风险评级","v":"89"},{"k":"风险授信","v":"61"},{"k":"总流程时长","v":"460"}]},{"Title":"流程干涉","Childs":[{"k":"","v":"无"}]}]}],"progress":[{"Name":"营销","Childs":[{"Name":"拜访","Percent":"6","Date":"2020/5/25"},{"Name":"确定方案","Percent":"10","Date":"2020/5/26"},{"Name":"收集资料","Percent":"15","Date":"2020/6/25"}]},{"Name":"评级","Childs":[{"Name":"撰写评级","Percent":"20","Date":"2020/7/25"},{"Name":"交支行","Percent":"25","Date":"2020/7/26"},{"Name":"退回","Percent":"21","Date":"2020/7/28"},{"Name":"交分行","Percent":"30","Date":"2020/10/25"},{"Name":"评级完成","Percent":"35","Date":"2020/11/25"}]},{"Name":"授信","Childs":[{"Name":"撰写授信","Percent":"40","Date":"2021/1/5"},{"Name":"交支行","Percent":"45","Date":"2021/1/15"},{"Name":"退回","Percent":"40","Date":"2021/1/25"},{"Name":"交分行","Percent":"50","Date":"2021/1/26"},{"Name":"上会","Percent":"55","Date":"2021/1/27"},{"Name":"授信批复","Percent":"60","Date":"2021/1/28"}]},{"Name":"提款","Childs":[{"Name":"撰写提款","Percent":"61","Date":"2021/2/28"},{"Name":"交支行","Percent":"70","Date":"2021/3/28"},{"Name":"退回","Percent":"62","Date":"2021/4/28"},{"Name":"提款批复","Percent":"78","Date":"2021/5/28"}]},{"Name":"放款","Childs":[{"Name":"提交","Percent":"80","Date":"2021/6/28"},{"Name":"退回","Percent":"81","Date":"2021/7/28"},{"Name":"放行","Percent":"100","Date":"2021/8/28"}]}]}`), &d )
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
	e = f.AddPicture("Sheet1", "A33", "C:/.repos/cyf-cloud.back/.vp/1.png", "")
	if e != nil {
		fmt.Println(e)
	}
	// 根据指定路径保存文件
	e = f.SaveAs("C:/.repos/cyf-cloud.back/.vp/1.xlsx")
	if e != nil {
		fmt.Println(e)
	}
}


