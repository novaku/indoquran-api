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
	IP             string `json:"ip"`
	ContinentCode  string `json:"continent_code"`
	ContinentName  string `json:"continent_name"`
	CountryCode2   string `json:"country_code2"`
	CountryCode3   string `json:"country_code3"`
	CountryName    string `json:"country_name"`
	CountryCapital string `json:"country_capital"`
	StateProv      string `json:"state_prov"`
	District       string `json:"district"`
	City           string `json:"city"`
	Zipcode        string `json:"zipcode"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	IsEu           bool   `json:"is_eu"`
	CallingCode    string `json:"calling_code"`
	CountryTld     string `json:"country_tld"`
	Languages      string `json:"languages"`
	CountryFlag    string `json:"country_flag"`
	GeonameID      string `json:"geoname_id"`
	Isp            string `json:"isp"`
	ConnectionType string `json:"connection_type"`
	Organization   string `json:"organization"`
	Currency       struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currency"`
	TimeZone struct {
		Name            string  `json:"name"`
		Offset          int     `json:"offset"`
		CurrentTime     string  `json:"current_time"`
		CurrentTimeUnix float64 `json:"current_time_unix"`
		IsDst           bool    `json:"is_dst"`
		DstSavings      int     `json:"dst_savings"`
	} `json:"time_zone"`
}
