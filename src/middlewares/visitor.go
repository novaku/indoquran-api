package middlewares

import (
	"fmt"
	"time"

	"bitbucket.org/indoquran-api/src/config"
	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/helpers"
	"bitbucket.org/indoquran-api/src/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
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
	cc := c.Copy()

	go func() {
		// get userid from session
		session := sessions.Default(cc)
		userID := session.Get(config.Config.Session.UserID)
		ipData, err := helpers.IPToCountry(cc, cc.ClientIP())
		if err != nil {
			mlog.Error(err)
		}

		// get body request
		visitor := models.Visitor{
			ID:        primitive.NewObjectID(),
			UserID:    fmt.Sprintf("%v", userID),
			IP:        cc.ClientIP(),
			IPData:    ipData,
			Path:      cc.Request.URL.Path,
			URL:       cc.Request.URL.String(),
			Method:    cc.Request.Method,
			UserAgent: cc.Request.Header.Get("User-Agent"),
			Ref:       cc.Request.Header.Get("Referer"),
			Time:      time.Now(),
		}

		visitorCollection.InsertOne(cc, visitor)
	}()

	c.Next()
}
