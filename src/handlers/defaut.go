package handlers

import (
	"log"

	"bitbucket.org/indoquran-api/src/config"
	"github.com/gin-gonic/gin"
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

// MongoConfig : Get DB from Mongo Config
func MongoConfig() *mongo.Database {
	db, err := config.GetMongoDB()
	if err != nil {
		log.Println("ERROR : ", err)
	}
	return db
}
