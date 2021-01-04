package user

import (
	"time"

	model_user "bitbucket.org/indoquran-api/src/models/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func findOrInsertDocument(c *gin.Context, userForm *formPost) (string, error) {
	user := model_user.User{}
	filter := bson.M{
		"email": bson.M{
			"$eq": userForm.Email, // check if email already registered
		},
	}

	usr := userCollection.FindOne(c, filter)
	if usr.Err() == nil {
		err := usr.Decode(&user)
		if err != nil {
			return "", err
		}

		loginHistory(c, user.ID)
		return user.ID.Hex(), nil
	}

	model := model_user.User{
		ID:         primitive.NewObjectID(),
		Name:       userForm.Name,
		Email:      userForm.Email,
		FacebookID: userForm.FacebookID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	data, err := userCollection.InsertOne(c, &model)
	if err != nil {
		return "", err
	}

	loginHistory(c, data.InsertedID.(primitive.ObjectID))
	return data.InsertedID.(primitive.ObjectID).Hex(), nil
}

func loginHistory(c *gin.Context, userID primitive.ObjectID) error {
	history := model_user.History{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := loginHistoryCollection.InsertOne(c, history)
	if err != nil {
		return err
	}

	return nil
}
