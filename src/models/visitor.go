package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Visitor : visitor structure
type Visitor struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	IP        string             `bson:"ip"`
	Path      string             `bson:"path"`
	URL       string             `bson:"url"`
	Method    string             `bson:"method"`
	UserAgent string             `bson:"user_agent"`
	Ref       string             `bson:"ref"`
	Time      time.Time          `bson:"time"`
}
