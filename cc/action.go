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
	"github.com/kpango/glg"
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
	ActionGroupFunc func( ActionGroup ) error
	ActionFunc func( ActionPackage ) ( HttpErrReturn, StatusCode )
)

var (
	postHandlers map[string] *ActionFunc
	getHandlers map[string] *ActionFunc

	actionGroupHandlers map[string] ActionGroupFunc
)

func init() {
	postHandlers = make( map[string] *ActionFunc )
	getHandlers = make( map[string] *ActionFunc )

	actionGroupHandlers = make(map[string]ActionGroupFunc)
}

// 添加一个业务逻辑组
// 所有的 action 将在 RegisterActions() 被调用时启用
func AddActionGroup( groupPath string, actionFunc ActionGroupFunc) {
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

// 添加一个Post请求
func ( a ActionGroup ) POST( path string, handler ActionFunc ) {
	glg.Log( "[action] POST: ", a.Path + path )
	http.HandleFunc( a.Path + path, mwh.WrapPost(
		func( w http.ResponseWriter, r *http.Request ) {
			her, status := handler( ActionPackage{ R: r, W: &w } )
			HttpReturnHER( &w, &her, status)
		} ) )
	postHandlers[path] = &handler
}

// 添加一个Get请求
func ( a ActionGroup ) GET( path string, handler ActionFunc ) {
	glg.Log( "[action] GET: ", a.Path + path )
	http.HandleFunc( a.Path + path, mwh.WrapGet(
		func( w http.ResponseWriter, r *http.Request ) {
			her, status := handler( ActionPackage{ R: r, W: &w } )
			HttpReturnHER( &w, &her, status)
		} ) )
	getHandlers[path] = &handler
}

func ( pap *ActionPackage )SetCookie( cookie *http.Cookie ) {
	http.SetCookie( *pap.W, cookie )
}

func ( pap *ActionPackage ) GetCookie( key string ) ( string, error ) {
	cl, e  := pap.R.Cookie("atk")
	if e != nil {
		glg.Error("atk not found. it may be a post proxy problem")
		return "", e
	}
	atk := cl.Value
	glg.Success("atk is (" + atk + ")")
	return atk, e
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