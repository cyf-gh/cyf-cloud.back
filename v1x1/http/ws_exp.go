package http

import (
	"../../cc"
	"github.com/gorilla/websocket"
	"github.com/kpango/glg"
)

func init() {
	cc.AddActionGroup( "/v1x1/ws/test", func( a cc.ActionGroup ) error {
		a.WS( "/echo", func( ap cc.ActionPackage, c *websocket.Conn ) ( e error ) {
			for {
				mt, message, err := c.ReadMessage()
				if err != nil {
					glg.Error( "WS READ MSG:", err )
					break
				}
				glg.Log("[" + ap.R.URL.Path + "] WS RECV: ", string(message) )

				err = c.WriteMessage( mt, message )
				if err != nil {
					glg.Error( "WS WRITE MSG:", err )
					break
				}
			}
			return nil
		} )
		return nil
	} )
}