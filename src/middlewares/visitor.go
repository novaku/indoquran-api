package middlewares

import (
	"fmt"
	"time"

	"bitbucket.org/indoquran-api/src/config"
	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	visitorCollName = "visitor"
)

var (
	db                = *handlers.MongoConfig()
	visitorCollection = db.Collection(visitorCollName)
)

// VisitorLogger : middleware to log visitor acccess
func VisitorLogger(c *gin.Context) {
	// get userid from session
	session := sessions.Default(c)
	userID := session.Get(config.Config.Session.UserID)

	// get body request
	visitor := models.Visitor{
		ID:        primitive.NewObjectID(),
		UserID:    fmt.Sprintf("%v", userID),
		IP:        c.ClientIP(),
		Path:      c.Request.URL.Path,
		URL:       c.Request.URL.String(),
		Method:    c.Request.Method,
		UserAgent: c.Request.Header.Get("User-Agent"),
		Ref:       c.Request.Header.Get("Referer"),
		Time:      time.Now(),
	}

	visitorCollection.InsertOne(c, visitor)

	c.Next()
}
