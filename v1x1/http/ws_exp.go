package http

import (
	"../../cc"
	"github.com/kpango/glg"
)

func init() {
	cc.AddActionGroup( "/v1x1/ws/test", func( a cc.ActionGroup ) error {
		a.WS( "/echo", func( ap cc.ActionPackage, aws cc.ActionPackageWS ) ( e error ) {
			c := aws.C
			for {
				mt, message, err := c.ReadMessage()
				if err != nil {
					glg.Error( "WS READ:", err )
					break
				}
				glg.Log("[" + ap.R.URL.Path + "] WS RECV: ", string(message) )

				err = c.WriteMessage( mt, message )
				if err != nil {
					glg.Error( "WS WRITE:", err )
					break
				}
			}
			return nil
		} )
		return nil
	} )
}