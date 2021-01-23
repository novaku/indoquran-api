package user

import (
	"net/http"

	"bitbucket.org/indoquran-api/src/config"
	"bitbucket.org/indoquran-api/src/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoginUser : Login user controller
func LoginUser(c *gin.Context) {
	form := formPost{}
	session := sessions.Default(c)

	if err := c.ShouldBind(&form); err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusBadRequest, "Error Get Body", err.Error())
		return
	}

	id, err := findOrInsertDocument(c, &form)
	if err != nil {
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed login user", err.Error())
		return
	}

	session.Set(config.Config.Session.UserID, id)
	session.Save()

	handlers.DefaultResponse(c, http.StatusOK, "Succes Login User", id)
}

// GetUser : get one user controller
func GetUser(c *gin.Context) {
	id := c.Param("id")
	idParse, errParse := primitive.ObjectIDFromHex(id)
	if errParse != nil {
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Parse Param", errParse.Error())
		return
	}

	res := userCollection.FindOne(c, bson.M{"_id": &idParse})
	if res.Err() != nil {
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed Get One User", res.Err().Error())
		return
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success Get One User", &res)
}

// LogoutUser : logout user
func LogoutUser(c *gin.Context) {
	form := formLogout{}
	session := sessions.Default(c)

	if err := c.ShouldBind(&form); err != nil {
		handlers.DefaultResponse(c, http.StatusBadRequest, "Bad request logout data", err.Error())
		return
	}

	err := findByUserID(c, form.UserID)

	if err != nil {
		handlers.DefaultResponse(c, http.StatusBadRequest, "Unable to find user in database", err.Error())
		return
	}

	session.Delete(config.Config.Session.UserID)
	session.Save()
	c.SetCookie(config.Config.Session.Name, "", -1, "/", c.Request.URL.Host, false, false)

	handlers.DefaultResponse(c, http.StatusOK, "Success", nil)
}
