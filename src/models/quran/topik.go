package quran

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Topik : struct of topik
type Topik struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	TopicID   int                `bson:"topicID" json:"topicID"`
	ParentID  int                `bson:"ayat" json:"parentID"`
	IsTitle   int                `bson:"isTitle" json:"isTitle"`
	Text      string             `bson:"text" json:"text"`
	Isi       string             `bson:"isi" json:"isi"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
