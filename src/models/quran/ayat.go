package quran

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ayat : struct of ayat
type Ayat struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Nomor     int                `bson:"nomor" json:"nomor"`
	Surat     int                `bson:"surat" json:"surat"`
	Juz       int                `bson:"juz" json:"juz"`
	TxtAR     string             `bson:"txt_ar" json:"txt_ar"`
	TxtID     string             `bson:"txt_id" json:"txt_id"`
	TxtIDT    string             `bson:"txt_idt" json:"txt_idt"`
	TxtTafsir string             `bson:"txt_tafsir" json:"txt_tafsir"`
	Image     string             `bson:"image" json:"image"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

//Ayats : collections of ayat
type Ayats []Ayat
