package handlers

import (
	"log"

	"bitbucket.org/indoquran-api/src/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// DefaultResponse : default success response
func DefaultResponse(c *gin.Context, httpCode int, message string, obj interface{}) {
	c.JSON(httpCode, gin.H{
		"message": message,
		"data":    obj,
	})
	return
}

// MongoInstance : Get DB from Mongo Config
func MongoInstance() *mongo.Database {
	db, err := config.GetMongoDB()
	if err != nil {
		log.Println("Connect to Mongo failed, ERROR : ", err)
	}
	return db
}

// RedisInstance : get cache from Redis config
func RedisInstance() *redis.Client {
	opt, err := config.GetRedis()
	if err != nil {
		log.Println("Connect to redis failed, ERROR : ", err)
	}

	return redis.NewClient(opt)
}
