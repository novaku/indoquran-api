package quran

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tafsir : struct of tafsir kemenag
type Tafsir struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Surat      int                `bson:"surat" json:"surat"`
	Ayat       int                `bson:"ayat" json:"ayat"`
	TextTafsir string             `bson:"text_tafsir" json:"text_tafsir"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// Tafsirs : collections of tafsir kemenag
type Tafsirs []Tafsir
