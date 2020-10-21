package run_mode

import (
	"../config"
	"../v1x1/cache"
)

func HttpInit() {
	if config.IsRunModeDev() {
		// 方便测试的令牌，cyfhaoshuai
		cache.Set( "cyfhaoshuaicyfhaoshuaicyfhaoshuaicyfhaoshuai", 1 )
	}
}
