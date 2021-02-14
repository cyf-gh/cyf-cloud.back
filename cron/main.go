package cron

import (
	"github.com/robfig/cron/v3"
)

func init() {
	c := cron.New()
	c.AddFunc("@every 0.5h", DumpTasks)
}
