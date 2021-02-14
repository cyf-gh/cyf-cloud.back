package cron

import (
	"../dm_1"
	"../v1x1/cache"
	"encoding/json"
	"github.com/kpango/glg"
)

const (
	KeyCronDumpTask = "cc_cron_dump_tasks"
)

func DumpTasks()  {
	bs, e := json.Marshal( dm_1.TaskSharedList ); if e != nil { glg.Error( "[cron] in dumpTasks", e ); return }
	cache.Set( KeyCronDumpTask, string( bs ) )
	glg.Info("[cron] dumpTasks finished")
}

