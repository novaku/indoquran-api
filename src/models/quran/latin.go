package quran

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ayat : struct of ayat
type Latin struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	AyatID int                `bson:"ayat_id" json:"ayat_id"`
	Surat  int                `bson:"surat" json:"surat"`
	Latin  string             `bson:"latin" json:"latin"`
}

//Ayats : collections of ayat
type Latins []Latin
