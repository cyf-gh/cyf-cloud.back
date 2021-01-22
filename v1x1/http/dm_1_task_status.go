package http

import (
	"../../cc"
	"../../dm_1"
)

func init() {
	// \brief 返回当前任务的情况
	// \sa dm_1.DMTaskStatus
	// \body see TaskInfo
	// 如果TaskInfo.TaskIndex为负数，则返回所有任务情况
	cc.AddActionGroup( "/v1x1/dm/1/tasks", func( a cc.ActionGroup ) error {
		a.WS( "/status/ws", func( ap cc.ActionPackage, aws cc.ActionPackageWS ) ( e error ) {
			type (
				TaskInfo struct {
					TaskName string
					TaskIndex int64
				}
			)
			for {
				ti := TaskInfo{}
				e = aws.ReadJson( &ti ); if e != nil { break }
				if ti.TaskIndex == -1 {
					e = aws.WriteJson( dm_1.TaskSharedList.Lists[ti.TaskName][ti.TaskIndex] ); if e != nil { break }
				} else {
					e = aws.WriteJson( dm_1.TaskSharedList.Lists ); if e != nil { break }
				}
			}
			return nil
		} )
		return nil
	} )
}