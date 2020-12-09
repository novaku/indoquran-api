package user

import (
	"net/http"
	"time"

	"bitbucket.org/indoquran-api/src/handlers"
	model_user "bitbucket.org/indoquran-api/src/models/user"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllUser : Get All User Endpoint
func GetAllUser(c *gin.Context) {
	res, err := userCollection.Find(c, bson.M{})

	if err != nil {
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed Get All User", err.Error())
		return
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success Get All User", &res)
}

// GetUser : Get User Endpoint
func GetUser(c *gin.Context) {
	id := c.Param("id")
	idParse, errParse := primitive.ObjectIDFromHex(id)
	if errParse != nil {
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Parse Param", errParse.Error())
		return
	}

	res := userCollection.FindOne(c, bson.M{"id": &idParse})

	handlers.DefaultResponse(c, http.StatusOK, "Success Get One User", &res)
}

// CreateUser : Create User Endpoint
func CreateUser(c *gin.Context) {
	user := model_user.User{}
	err := c.Bind(&user)

	if err != nil {
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Get Body", err.Error())
		return
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	res, err := userCollection.InsertOne(c, user)

	if err != nil {
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Insert User", err.Error())
		return
	}

	handlers.DefaultResponse(c, http.StatusOK, "Succes Insert User", &res)
}

// UpdateUser : Update User Endpoint
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Parse Param", err.Error())
		return
	}

	user := model_user.User{}
	err = c.Bind(&user)

	if err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Bind Param", err.Error())
		return
	}

	user.UpdatedAt = time.Now()

	res := userCollection.FindOneAndUpdate(c, bson.M{"_id": idObject}, user)

	handlers.DefaultResponse(c, http.StatusOK, "Succes Update User", &res)
}

// DeleteUser : Delete User Endpoint
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Parse Param", err.Error())
		return
	}

	res, err := userCollection.DeleteOne(c, bson.M{"_id": idObject})

	if err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Error Delete User", err.Error())
		return
	}

	handlers.DefaultResponse(c, http.StatusOK, "Succes Delete User", &res)
}
