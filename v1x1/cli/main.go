package cli

import (
	c "../../cli"
	cfg "../../config"
	"../cache"
)

func Init() {
	c.Register( "rget", &c.CliFuncPack{ F: GetRedisKeyValue, Desc: "Get value by key in redis", Group: "redis" })
	c.Register("mode", &c.CliFuncPack{
		F:     ChangeMode,
		Desc:  "Change server mode",
		Group: "server",
	})
}

func GetRedisKeyValue( ts []string ) error {
	v, e := cache.Get( ts[0] )
	println( v )
	return e
}

func ChangeMode( ts []string ) error {
	mode := ts[0]
	cfg.RunMode = mode
	println( "mode switched to \"" + mode + "\"")
	return nil
}