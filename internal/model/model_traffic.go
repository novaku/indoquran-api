package model

import (
	"gorm.io/gorm"
)

type ApiTrafficLog struct {
	gorm.Model
	IPAddress      string `gorm:"not null"`
	Endpoint       string `gorm:"not null"`
	Duration       string
	HTTPMethod     string `gorm:"enum('GET','POST','PUT','DELETE','PATCH','OPTIONS','HEAD');not null"`
	RequestPayload string
	ResponseStatus int `gorm:"not null"`
	ResponseBody   string
	UserAgent      string
	Referer        string
}

func (ApiTrafficLog) TableName() string {
	return "api_traffic_log"
}
