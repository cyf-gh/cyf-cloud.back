package cli

import (
	c "../../cli"
	"../cache"
)

func Init() {
	c.Register( "rget", &c.CliFuncPack{ GetRedisKeyValue, "Get value by key in redis", "redis" })
}

func GetRedisKeyValue( ts []string ) error {
	v, e := cache.Get( ts[0] )
	println( v )
	return e
}
