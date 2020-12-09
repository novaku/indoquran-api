package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User : User structure
type User struct {
	ID        primitive.ObjectID `bson:"id"`
	Name      string             `bson:"name"`
	Address   string             `bson:"address"`
	Age       int                `bson:"age"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

//Users : list of users
type Users []User
