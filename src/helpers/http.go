package helpers

import (
	"bufio"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"bitbucket.org/indoquran-api/src/config"
	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/models"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
)

const (
	redisKeyIPToCountry = "redis:iptocuntry:"
)

var cache = *handlers.RedisConfig()

// IPToCountry : IP to country data
func IPToCountry(c *gin.Context, ip string) (*models.IPToCountryStruct, error) {
	client := &http.Client{}
	ipData := models.IPToCountryStruct{}

	val, err := cache.Get(c, redisKeyIPToCountry+ip).Result()
	if err != nil || val == "" {
		request, err := http.NewRequest("GET", "https://api.ipgeolocation.io/ipgeo?apiKey="+config.Config.Secret.Geolication+"&ip="+ip, nil)
		if err != nil {
			mlog.Error(err)
			return nil, err
		}

		request.Header.Set("User-Agent", randomUserAgents())

		response, err := client.Do(request)
		if err != nil {
			mlog.Error(err)
		}

		defer response.Body.Close()

		err = json.NewDecoder(response.Body).Decode(&ipData)
		if err != nil {
			mlog.Error(err)
			return nil, err
		}

		mlog.Info("JSON response from IP to country: %+v", ipData)

		b, err := json.Marshal(&ipData)
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

	if val != "" {
		mlog.Info("GET redis key : %s", redisKeyIPToCountry+ip)

		err = json.Unmarshal([]byte(val), &ipData)
		if err != nil {
			panic(err)
		}
	}

	return &ipData, nil
}

func randomUserAgents() string {
	f, err := os.Open("resources/user-agents.txt")
	agentCount := 0
	agents := []string{}

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		agents = append(agents, scanner.Text())
		agentCount++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return agents[rand.Intn(len(agents)-1)]
}
