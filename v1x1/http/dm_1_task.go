package http

import (
	"../../cc"
	"../../cc/err"
	"../../dm_1"
	"strconv"
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/task", func( a cc.ActionGroup ) error {
		// \brief 返回当前任务的情况
		// \sa dm_1.DMTaskStatus
		// \arg[TaskInfo] see TaskInfo 如果TaskInfo.Index为-1，则返回所有任务情况
		// \note 无特殊情况都使用-1的参数
		// \return []DMForableTaskSharedList
		a.WS( "/status/ws", func( ap cc.ActionPackage, aws cc.ActionPackageWS ) ( e error ) {
			type (
				TaskInfo struct {
					Name string
					Index int64
				}
			)
			for {
				ti := TaskInfo{}
				e = aws.ReadJson( &ti ); if e != nil { break }
				if ti.Index != -1 {
					e = aws.WriteJson( dm_1.TaskSharedList.Lists[ti.Name][ti.Index] ); if e != nil { break }
				} else {
					e = aws.WriteJson( dm_1.TaskSharedList.ToForable() ); if e != nil { break }
				}
			}
			return nil
		} )
		// \brief 终止某个任务
		// \arg[name] 任务名字
		// \arg[i] 任务索引位置
		// \return ok
		a.GET( "/abort", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			name := ap.GetFormValue("name")
			index := ap.GetFormValue( "i" )
			i, e := strconv.Atoi( index ); err.Assert( e )
			dm_1.TaskSharedList.Lists[name][i].Terminal = true
			return cc.HerOk()
		} )
		// \brief 暂停某个任务
		// \note 任务能否暂停请阅读创建任务的API注释；有的任务不可暂停；任务通常用自旋而非阻塞进行暂停自旋间隔请阅读创建任务的API注释
		// 重复暂停 NSF
		// \arg[name] 任务名字
		// \arg[i] 任务索引位置
		// \return ok
		a.GET( "/pause", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			name := ap.GetFormValue("name")
			index := ap.GetFormValue( "i" )
			i, e := strconv.Atoi( index ); err.Assert( e )
			dm_1.TaskSharedList.Lists[name][i].Pause = true
			return cc.HerOk()
		} )
		// \brief 恢复某个任务
		// \note 重复恢复 NSF
		// \arg[name] 任务名字
		// \arg[i] 任务索引位置
		// \return ok
		a.GET( "/resume", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			name := ap.GetFormValue("name")
			index := ap.GetFormValue( "i" )
			i, e := strconv.Atoi( index ); err.Assert( e )
			dm_1.TaskSharedList.Lists[name][i].Pause = false
			return cc.HerOk()
		} )
		return nil
	} )
}