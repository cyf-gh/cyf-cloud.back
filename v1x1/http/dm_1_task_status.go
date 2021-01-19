package http

import (
	"../../cc"
)

type (
	TaskInfo struct {
		TaskName string
		TaskIndex int64
	}
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/task", func( a cc.ActionGroup ) error {
		a.WS( "/status", func( ap cc.ActionPackage, aws cc.ActionPackageWS ) ( e error ) {
			for {
				ti := TaskInfo{}
				e = aws.ReadJson( &ti ); if e != nil { break }


			}
			return nil
		} )
		return nil
	} )
}