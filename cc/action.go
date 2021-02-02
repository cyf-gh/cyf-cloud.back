package cc

/**

Example:

//  func init() {
//
//  cc.AddActionGroup( "/foo", func( a cc.ActionGroup ) error {
//		a.GET( "/aaa/bbb", func( w http.ResponseWriter, r *http.Request ) {
// 			// ...业务逻辑...
//		}
//		a.POST( "/aaa/bbb", func( w http.ResponseWriter, r *http.Request ) {
// 			// ...业务逻辑...
//		}
// 	}
//	return nil
//
//  }

}
 */

import (
	mwh "../middleware/helper"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
)

type (
	ActionGroup struct {
		Path string
	}
	ActionPackage struct {
		R *http.Request
		W *http.ResponseWriter
	}
	ActionPackageWS struct {
		C *websocket.Conn
	}
	ActionGroupFunc func( ActionGroup ) error
	ActionFunc func( ActionPackage ) ( HttpErrReturn, StatusCode )
	ActionFuncWS func( ActionPackage, ActionPackageWS ) error
)

var (
	postHandlers map[string] *ActionFunc
	getHandlers map[string] *ActionFunc
	wsHandlers map[string] *ActionFuncWS

	actionGroupHandlers map[string] ActionGroupFunc
)

func init() {
	postHandlers = make( map[string] *ActionFunc )
	getHandlers = make( map[string] *ActionFunc )
	wsHandlers = make( map[string] *ActionFuncWS )

	actionGroupHandlers = make(map[string]ActionGroupFunc)
}

func ( R ActionPackage ) GetFormValue( key string ) string {
	return R.R.FormValue(key)
}

func ( R ActionPackage ) GetBodyUnmarshal( v interface{} ) error {
	b, e := ioutil.ReadAll( R.R.Body ); if e != nil { return e }
	e = json.Unmarshal( b, v ); if e != nil { return e }
	return nil
}

// 添加一个业务逻辑组
// 所有的 action 将在 RegisterActions() 被调用时启用
func AddActionGroup( groupPath string, actionFunc ActionGroupFunc) {
	checkPathWarning( groupPath )
	if _, ok := actionGroupHandlers[groupPath]; ok {
		glg.Warn("action group:", groupPath, "already exists, recovered.")
	}
	actionGroupHandlers[groupPath] = actionFunc
}

// 启用所有路由
func RegisterActions() error {
	for k, a := range actionGroupHandlers {
		if e := a( ActionGroup{Path: k} ); e != nil {
			glg.Error("in action group:", k)
			return e
		}
	}
	return nil
}

// cc标准的路径都均为 开头 /xxx 或 空
// 路径最后一个字符不得为 /
// 检查不符合标准仅在输出警告
func checkPathWarning( path string ) {
	if path == "" {
		return
	}
	if path[:1] != "/" || path[len(path)-1:] == "/" {
		glg.Warn( "url: ", path, " may not correct; are you sure it was the expected url path?")
	}
}

// 添加一个Post请求
func ( a ActionGroup ) POST( path string, handler ActionFunc ) {
	checkPathWarning( path )
	glg.Log( "[action] POST: ", a.Path + path )
	http.HandleFunc( a.Path + path, mwh.WrapPost(
		func( w http.ResponseWriter, r *http.Request ) {
			her, status := handler( ActionPackage{ R: r, W: &w } )
			HttpReturnHER( &w, &her, status, r.URL.Path )
		} ) )
	postHandlers[path] = &handler
}

// 添加一个Get请求
func ( a ActionGroup ) GET( path string, handler ActionFunc ) {
	checkPathWarning( path )
	glg.Log( "[action] GET: ", a.Path + path )
	http.HandleFunc( a.Path + path, mwh.WrapGet(
		func( w http.ResponseWriter, r *http.Request ) {
			her, status := handler( ActionPackage{ R: r, W: &w } )
			HttpReturnHER( &w, &her, status, r.URL.Path	)
		} ) )
	getHandlers[path] = &handler
}

// 添加一个websocket请求
// cc规范：必须在请求路径末端添加ws字段来提示这一请求为websocket请求
// 例：/imai_mami/no/koto/ga/suki/ws
func ( a ActionGroup ) WS( path string, handler ActionFuncWS ) {
	checkPathWarning( path )
	glg.Log( "[action] WS: ", a.Path + path )
	http.HandleFunc( a.Path + path, mwh.WrapWS( func( w http.ResponseWriter, r *http.Request ) {
			glg.Log("["+ a.Path + path  +"] "+ "WS: START UPGRADE")

			ug := websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				// 交给 nginx 处理源问题
				CheckOrigin: func( r *http.Request ) bool {
					return true
				},
			}
			c, e := ug.Upgrade( w, r, nil )
			if e != nil { glg.Error("["+ a.Path + path  +"] " + "WS UPGRADE: ", e); return }
			defer c.Close()

			if e = handler( ActionPackage{ R: r, W: &w }, ActionPackageWS{C:c} ); e != nil { glg.Error( e ) }
			glg.Info( "["+ a.Path + path  +"] " + "WS CLOSED" )
		} ) )
	wsHandlers[path] = &handler
}

func resp(w* http.ResponseWriter, msg string) {
	(*w).Write([]byte(msg))
}

// 只返回data，不返回其他的任何信息
// DO: DATA ONLY
func ( a ActionGroup ) GET_DO( path string, handler ActionFunc ) {
	glg.Log( "[action] GET_DO: ", a.Path + path )
	http.HandleFunc( a.Path + path, mwh.WrapGet(
		func( w http.ResponseWriter, r *http.Request ) {
			her, _ := handler( ActionPackage{ R: r, W: &w } )
			resp( &w, her.Data )
		} ) )
	getHandlers[path] = &handler
}

func ( pap *ActionPackage )SetCookie( cookie *http.Cookie ) {
	http.SetCookie( *pap.W, cookie )
}

func ( pap *ActionPackage ) GetCookie( key string ) ( string, error ) {
	cl, e  := pap.R.Cookie( key )
	if e != nil {
		glg.Error("key not found. it may be a post proxy problem")
		return "", e
	}
	res := cl.Value
	glg.Success("COOKIE ["+key+"] : (" + res + ")")
	return res, e
}

// 将ws读取数据转化为json
// error总是断连错误
func ( pR *ActionPackageWS ) ReadJson( v interface{} ) ( e error ) {
	mt, b, e := pR.C.ReadMessage(); if e != nil { return e }
	switch mt {
	case websocket.BinaryMessage:
		glg.Warn("WS: reading binary message but try to unmarshal it")
	case websocket.CloseMessage:
		glg.Log("WS closed")
		return errors.New("WS closed")
	}
	e = json.Unmarshal( b, v ); if e != nil { return e }
	return nil
}

// 将ws的读取数据转化为字符串
func ( pR *ActionPackageWS ) ReadString() ( string, error ) {
	mt, b, e := pR.C.ReadMessage(); if e != nil { return "", e }
	switch mt {
	case websocket.BinaryMessage:
		glg.Warn("WS: reading binary message but try to stringify it")
	case websocket.CloseMessage:
		glg.Log("WS closed")
		return "", errors.New("WS closed")
	}
	return string( b ), nil
}

func ( pR *ActionPackageWS ) ReadBinary() ( []byte, error ) {
	mt, b, e := pR.C.ReadMessage(); if e != nil { return nil, e }
	switch mt {
	case websocket.CloseMessage:
		glg.Log("WS closed")
		return nil, errors.New("WS closed")
	}
	return b, nil
}

func ( pR *ActionPackageWS ) WriteJson( data interface{} ) ( e error ) {
	jn, e := json.Marshal( data );  if e != nil { return e }
	e = pR.C.WriteMessage( websocket.TextMessage, jn ) ;  if e != nil { return e }
	return nil
}

func ( pR *ActionPackageWS ) WriteString( str string ) ( e error ) {
	e = pR.C.WriteMessage( websocket.TextMessage, []byte(str) ) ;  if e != nil { return e }
	return nil
}

func ( pR *ActionPackageWS ) WriteBinary( b []byte ) ( e error ) {
	e = pR.C.WriteMessage( websocket.BinaryMessage, b ) ;  if e != nil { return e }
	return nil
}

/*
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// 注册所有的Action
func RegisterAction( action *ActionGroup ) {
	aRf := reflect.ValueOf(&action).Elem()
	glg.Success( "[cc action] method number of action is ", aRf.NumMethod() )
	for i := 0; i < aRf.NumMethod(); i++  {
		glg.Info( "calling...", getFunctionName( aRf.Method( i ) ) )
		aRf.Method(i).Call(nil) // 不应该有参数
	}
}
*/