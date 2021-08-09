package http

import (
	"../../cc"
	"../../cc/err"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

func init() {
	cc.AddActionGroup( "/v1x1/practice", func( a cc.ActionGroup ) error {
		// \brief 返回题库
		// \type GET
		// \arg[dbname] database name 数据库文件名
		// \arg[stname] sheet name 表名
		// \arg[sr] start row 起始行，该行应当具有表头信息
		// \arg[head] 题目开始下标，可为缺省
		// \arg[end]  题目结束下标，可为缺省
		a.GET( "/db",  func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			dbname := ap.GetFormValue("dbname")
			stname := ap.GetFormValue("stname")
			startRow := ap.GetFormValue("sr" )
			sr, e := strconv.Atoi( startRow ); err.Assert( e )
			f, e := excelize.OpenFile( "./.practice/"+dbname ); err.Assert( e )

			rows, e := f.GetRows(stname); err.Assert( e )
			ds := make([]map[string]string, 0)

			for i, row := range rows {
				if i > sr {
					d := make(map[string]string)
					for j, col := range row {
						for x := 0; x < sr; x++ {
							if j >= len(rows[x]) {
								break
							}
							if rows[x][j] != "" {
								d[rows[x][j]] = col
							}
						}
					}
					if len(d) != 0 {
						ds = append(ds, d)
					}
				}
			}

			head := ap.GetFormValue("head")
			end := ap.GetFormValue("end")
			if head == "" {
				h, e := strconv.Atoi( head ); err.Assert( e )
				return cc.HerOkWithData( ds[h:] )
			} else if end == "" {
				en, e := strconv.Atoi( end ); err.Assert( e )
				return cc.HerOkWithData( ds[:en] )
			} else if head == "" && end =="" {
				return cc.HerOkWithData( ds )
			} else {
				h, e := strconv.Atoi( head ); err.Assert( e )
				en, e := strconv.Atoi( end ); err.Assert( e )
				return cc.HerOkWithData( ds[h:en] )
			}
		} )
		a.GET( "/db/sheets",  func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			dbname := ap.GetFormValue("dbname")
			f, e := excelize.OpenFile( "./.practice/"+dbname ); err.Assert( e )
			return cc.HerOkWithData( f.GetSheetList() )
		} )
		return nil
	} )
}
