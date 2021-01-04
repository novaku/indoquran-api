package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// History : model for login history
type History struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
