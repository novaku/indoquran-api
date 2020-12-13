package quran

import "go.mongodb.org/mongo-driver/bson/primitive"

// KataBijak : structure for kata_bijak  collection
type KataBijak struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Kata string             `bson:"kata" json:"kata"`
}
