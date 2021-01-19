package dm_1

import (
	cfg "../config"
	"errors"
	"github.com/kpango/glg"
	"time"
)

type (
	DMTaskStatus struct {
		TaskName string
		IsUni bool

		Progress int
		ProgressMax int
		ProgressStage string

		IsFinished bool
		StartTime time.Time
		TimeoutSec float64
		CurrentMsg string // 当前工作的信息
		Errors []error
	}
	DMTaskStatusList struct {
		Lists map[string] []DMTaskStatus
		Errors []error
	}
)

var (
	// http 请求-读写
	// websocket 只读
	TaskSharedList DMTaskStatusList
)

func init() {
	TaskSharedList = DMTaskStatusList{ Lists: map[string][]DMTaskStatus{} }
}

func ( pR *DMTaskStatus ) Error( e error ) {
	if e != nil {
		pR.Errors = append(pR.Errors, e)
	}
}

func ( pR *DMTaskStatus) Abort( e error ) {
	pR.Error( e )
	pR.Finished()
}

func ( pR *DMTaskStatus ) Finished() {
	pR.IsFinished = true
}

func ( R DMTaskStatus) UsedTime() time.Duration {
	return time.Since( R.StartTime )
}

func ( R DMTaskStatus) IsTimeout() bool {
	return R.UsedTime().Seconds() >= R.TimeoutSec
}

func ( R DMTaskStatus ) Msg( msg string ) {
	if cfg.IsRunModeDev() {
		glg.Log( R.TaskName + "->" + msg )
	}
	R.CurrentMsg = msg
}

func ( R DMTaskStatus ) MsgStage( msg string ) {
	if cfg.IsRunModeDev() {
		glg.Log( R.TaskName + "->" + msg )
	}
	R.ProgressStage = msg
}

// 尝试创建一个unique的任务时返回错误
func ( pR *DMTaskStatusList ) AddTask( taskName string, isUni bool, timeout float64, progressMax int ) error {
	if isUni {
		if _, exists := pR.Lists[taskName]; exists {
			e := errors.New("try to create a duplicated unique task:" + taskName)
			pR.Errors = append(pR.Errors, e)
			return e
		}
	}
	if _, exists := pR.Lists[taskName]; !exists {
		pR.Lists[taskName] = []DMTaskStatus{}
	}
	pR.Lists[taskName] = append( pR.Lists[taskName], DMTaskStatus{
		TaskName: 	taskName,
		IsUni:      isUni,
		Progress:   0,
		IsFinished: false,
		StartTime:  time.Now(),
		TimeoutSec: timeout,
		ProgressMax: progressMax,
	} )
	return nil
}
