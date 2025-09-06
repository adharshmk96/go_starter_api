package infra

import (
	"servicehub_api/pkg/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitGormDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&domain.Account{})
	db.AutoMigrate(&domain.AccountActivity{})

	return db
}
