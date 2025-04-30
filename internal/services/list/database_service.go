package list

import (
	"indoquran-api/pkg/database"

	"gorm.io/gorm"
)

type gormDatabaseService struct {
	db *gorm.DB
}

func NewGormDatabaseService() DatabaseService {
	return &gormDatabaseService{
		db: database.GetDB(),
	}
}
