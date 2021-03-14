// 防止服务器API访问过载中间件
package middlewareUtil

import (
	"sync"
	"time"
)

type (
	Record struct {
		Mutex sync.Mutex
		Q []time.Time
	}
)

var	(
	TGActiveRecorder map[string]map[string]Record   // uniType->url->{ mutex, []time }
													// 通过的请求时间会被存放于此
	TGRefusedRecorder map[string]map[string]Record 	// 未通过的请求会被存放于此
)

const (
	TGWaitPending = 1
	TGWaitReturnNow = 0
)

func init() {
	// 初始化代码最左层 map[int64]map[string]Record -> map[int64]
	TGActiveRecorder = map[string]map[string]Record{}
}

func getRecord( uniType, url string, recorder *map[string]map[string]Record ) *Record {
	if rec, ok := TGActiveRecorder[uniType]; !ok {
		// 无该访问记录，则添加
		(*recorder)[uniType] = map[string]Record{}
		(*recorder)[uniType][url] = Record{
			Mutex: sync.Mutex{},
			Q:     []time.Time{},
		}
		rr := rec[url]
		return &rr
	} else {
		rr := rec[url]
		return &rr
	}
}

// uniType 唯一标识符的一种，可为ip，用户，或global
// 返回true则为未过载 false为过载不应当继续运行
func TGIPRecordAccess( uniType, url string, prepFreq int64 ) ( _continue bool ) {
	r := getRecord( uniType, url, &TGActiveRecorder )

	r.Mutex.Lock()
	dura := r.Q[len(r.Q)-1].Sub(r.Q[0]).Seconds()
	nowFreq := int64( float64( len(r.Q) ) / dura )

	// 当只有一个时必能通过
	if len( r.Q ) == 1 || dura == 0 {
		r.Q = append( r.Q, time.Now() )
		goto UNLOCK
	}
	if nowFreq < prepFreq { // 未过载 添加至
		r.Q = append( r.Q, time.Now() )
	} else { // 过载
		rr := getRecord( uniType, url, &TGRefusedRecorder )
		rr.Q = append( rr.Q, time.Now() )
		_continue = false
	}
	UNLOCK:
	r.Mutex.Unlock()

	return
}

func CheckRecordAccessIP( max int64 ) {

}