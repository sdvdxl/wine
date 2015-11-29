package cache

import (
	"github.com/gin-gonic/contrib/cache"
	"github.com/sdvdxl/wine/sources/util/constant"
	"github.com/sdvdxl/wine/sources/util/log"
	"time"
)

var (
	Cache cache.CacheStore
)

func init() {
	log.Logger.Info("init cache ...")
	//	Cache = cache.NewRedisCache("redis-31pzn.q1.tenxcloud.net:48181", "Isdvdxl8", constant.LOGIN_EXPIRED_TIME)
	Cache = cache.NewRedisCache("54.223.168.107:61993", "Isdvdxl8", constant.LOGIN_EXPIRED_TIME)
	if err := Cache.Set("_test_is_valid_x_", "x", time.Second); err != nil {
		log.Logger.Error("error redis status,%v, will use memcached", err)
		Cache = cache.NewMemcachedStore([]string{}, constant.LOGIN_EXPIRED_TIME)
		if err = Cache.Set("_test_is_valid_x_", "x", time.Second); err != nil {
			log.Logger.Error("error redis status,%v, will use memry cache", err)
			Cache = cache.NewInMemoryStore(constant.LOGIN_EXPIRED_TIME)
		}
	}

	log.Logger.Info("cache inited")
}
