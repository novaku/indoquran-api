package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User : User structure
type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Name       string             `bson:"name"`
	Email      string             `bson:"email"`
	FacebookID string             `bson:"facebook_id"`
	Role       string             `bson:"role"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

//Users : list of users
type Users []User
