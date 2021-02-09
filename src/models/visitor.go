package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Visitor : visitor structure
type Visitor struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    string             `bson:"user_id"`
	IP        string             `bson:"ip"`
	IPData    *IPToCountryStruct `bson:"ipdata"`
	Path      string             `bson:"path"`
	URL       string             `bson:"url"`
	Method    string             `bson:"method"`
	UserAgent string             `bson:"user_agent"`
	Ref       string             `bson:"ref"`
	Time      time.Time          `bson:"time"`
}

// IPToCountryStruct : structure for ip to country data
type IPToCountryStruct struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}
