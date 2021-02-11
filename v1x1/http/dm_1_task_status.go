package http

import (
	"../../cc"
	"../../dm_1"
)

func init() {
	// \brief 返回当前任务的情况
	// \sa dm_1.DMTaskStatus
	// \arg[TaskInfo] see TaskInfo 如果TaskInfo.Index为-1，则返回所有任务情况
	// \note 无特殊情况都使用-1的参数
	// \return []DMForableTaskSharedList
	cc.AddActionGroup( "/v1x1/dm/1/tasks", func( a cc.ActionGroup ) error {
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
		return nil
	} )
}