package middleware

import (
	"github.com/kpango/glg"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
)

func init() {
}

/*
中间件例子：

func Foo( ...foos Foo)  middleware.MiddewareFunc  {
	return func( f http.HandlerFunc ) http.HandlerFunc {
		return func( w http.ResponseWriter, r *http.Request ) {
			//... before
			f(w, r)
			//... after
		}
	}
}
 */
type (
	MiddewareFunc func(http.HandlerFunc) http.HandlerFunc
)

var (
	mwFuncs [] MiddewareFunc
)

func getFuncName( f MiddewareFunc ) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return name
}

// 注册中间件
//
// 注册的中间件将线性添加
// 越是后添加的中间件越包裹在内层
func Register( f MiddewareFunc ) {
	mwFuncs = append(mwFuncs, f )
	glg.Log( "[Middleware]" + getFuncName( f ) + " registered(prefix)\t[index]:" + strconv.Itoa( len(mwFuncs ) - 1 ) )
}

// 注销中间件
func Unregister( f MiddewareFunc ) {
	for i, fc := range mwFuncs {
		if getFuncName( fc ) ==  getFuncName( f ) {
			mwFuncs = append(mwFuncs[:i], mwFuncs[i+1:]...)
			glg.Log( "[Middleware]" + getFuncName( f ) + " unregistered(prefix) [index]:" + strconv.Itoa( i ) )
			return
		}
	}
	glg.Log( "[Middleware]" + getFuncName( f ) + " tried unregister(prefix) but failed" )
}

// 使用所有中间件
//
// mws：为nil则不运行额外的中间件
// 总是包裹在最内层
func HandlerWrapFully( handler http.HandlerFunc, mws ...MiddewareFunc ) http.HandlerFunc {
	for _, m := range mwFuncs {
		handler = m( handler )
	}
	if mws != nil {
		for _, m := range mws {
			handler = m( handler )
		}
	}
	return handler
}

// 选用指定的中间件 ，顺序可自定义
//
// index：
// { 1, 3, 5 } 表示只调用第2个，第4个，第5个中间件，线性执行
//
// mws：为nil则不运行额外的中间件
// 总是包裹在最内层
func HandlerWrap( handler http.HandlerFunc, index []int, mws ...MiddewareFunc ) http.HandlerFunc {
	for _, i := range index {
		handler = mwFuncs[i]( handler )
	}
	if mws != nil {
		for _, m := range mws {
			handler = m( handler )
		}
	}
	return handler
}