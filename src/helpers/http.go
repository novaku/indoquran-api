package helpers

import (
	"net/http"
	"time"

	"bitbucket.org/indoquran-api/src/config"
	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/models"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"github.com/vmihailenco/msgpack"
)

const (
	redisKeyIPToCountry = "redis:iptocuntry:"
)

var cache = *handlers.RedisConfig()

// IPToCountry : IP to country data
func IPToCountry(c *gin.Context, ip string) (*models.IPToCountryStruct, error) {
	ipData := models.IPToCountryStruct{}
	val, err := cache.Get(c, redisKeyIPToCountry+ip).Result()
	if err != nil || val == "" {
		response, err := http.Get("https://demo.ip-api.com/json/" + ip)
		if err != nil {
			mlog.Error(err)
			return nil, err
		}

		defer response.Body.Close()

		msgpack.NewDecoder(response.Body).Decode(&ipData)

		b, err := msgpack.Marshal(&ipData)
		if err != nil {
			mlog.Error(err)
			return nil, err
		}

		ttl := time.Duration(config.Config.Cache.TTL) * time.Hour
		set, err := cache.SetNX(c, redisKeyIPToCountry+ip, string(b), ttl).Result()
		if !set && err != nil {
			mlog.Error(err)
			return nil, err
		}
	}

	return &ipData, nil
}
