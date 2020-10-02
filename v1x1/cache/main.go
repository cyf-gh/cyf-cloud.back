package cache

import (
	"../../config"
	"github.com/gomodule/redigo/redis"
	"github.com/kpango/glg"
	"os"
	"strconv"
	"time"
)

var RedisPool *redis.Pool

func Init( rc config.RedisConfig ) {
	RedisPool = &redis.Pool{
		MaxIdle: rc.MaxIdle,
		MaxActive: rc.MaxActive,
		Dial: func() (redis.Conn, error) {
			conn, e := redis.Dial("tcp", rc.Addr)
			if e != nil {
				glg.Error("fail initializing the redis pool: %s", e.Error() )
				os.Exit(1)
			}
			return conn, e
		},
	}
	glg.Success("redis pool creation")

	for i := 0; i < 5; i++ {
		_, e := RedisPool.Get().Do("PING")
		if e == nil {
			glg.Success( "redis ping" )
			return
		}
		time.Sleep( 1000 )
		glg.Warn("redis ping failed, retry after 1 second..")
		glg.Warn("retry time" + strconv.Itoa(i))
	}
	glg.Fail("redis ping")
	os.Exit(1)
}

func Set( key, value string ) (interface{}, error)  {
	return RedisPool.Get().Do("SET", key, value)
}

func Get( key string ) ( string, error) {
	return redis.String( RedisPool.Get().Do("GET", key ) )
}