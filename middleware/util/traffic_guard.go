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
	TGActiveRecorder map[string]map[string]*Record   // uniType->url->{ mutex, []time }
													// 通过的请求时间会被存放于此
	TGRefusedRecorder map[string]map[string]*Record 	// 未通过的请求会被存放于此

	TGMutex sync.Mutex
)

const (
	TGWaitPending = 1
	TGWaitReturnNow = 0
)

func init() {
	// 初始化代码最左层 map[int64]map[string]Record -> map[int64]
	TGActiveRecorder = map[string]map[string]*Record{}
	TGMutex = sync.Mutex{}
}

// map 非线程安全 不允许concurrent读写
func getRecord( uniType, url string ) *Record {
	if _, ok := TGActiveRecorder[uniType]; !ok {
		// 无uniType访问记录访问记录，则添加
		TGActiveRecorder[uniType] = map[string]*Record{}

	}
	if _, ok := TGActiveRecorder[uniType][url]; !ok {
		// 无url访问记录，则添加
		TGActiveRecorder[uniType][url] = &Record{
			Mutex: sync.Mutex{},
			Q:     []time.Time{},
		}
	}
	rr := TGActiveRecorder[uniType][url]
	return rr
}

// uniType 唯一标识符的一种，可为ip，user，或global
// 返回true则为未过载 false为过载不应当继续运行
// Freq 单位为秒
// 本算法无需出栈，只需要按时释放time队列即可
func TGRecordAccess( uniType, url string, prepFreq float64 ) ( nowFreq float64, _continue bool ) {
	TGMutex.Lock()
	_continue = true
	r := getRecord( uniType, url )
	TGMutex.Unlock()

	r.Mutex.Lock()
	var (
		dura int64
	)
	// 当只有一个时必能通过
	if len( r.Q ) == 0 || len( r.Q ) == 1 {
		r.Q = append( r.Q, time.Now() )
		goto UNLOCK
	}
	dura = 9999
	// 将区间内的时间控制在一秒内
	for dura >= 1000 && len( r.Q ) > 1 {
		dura = r.Q[len(r.Q)-1].Sub(r.Q[0]).Milliseconds()
		r.Q = append( r.Q[:0], r.Q[1:]... )
	}
	nowFreq = float64( len(r.Q) ) / float64( dura ) * 1000

	r.Q = append( r.Q, time.Now() )
	if nowFreq >= prepFreq {
		_continue = false
	}
	UNLOCK:
	r.Mutex.Unlock()

	return
}

func CheckRecordAccessIP( max int64 ) {

}