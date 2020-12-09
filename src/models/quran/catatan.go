package quran

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Catatan : struct of catatan DEPAG
type Catatan struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Nomor       int                `bson:"nomor" json:"nomor"`
	Surat       int                `bson:"surat" json:"surat"`
	Ayat        int                `bson:"ayat" json:"ayat"`
	TeksCatatan string             `bson:"teks_catatan" json:"teks_catatan"`
	TeksQuran   string             `bson:"teks_quran" json:"teks_quran"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

//Catatans : collections of catatan
type Catatans []Catatan
