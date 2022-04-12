package database

import (
	"github.com/faunX/ogame/cmd/ogamed/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
	})

	db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Bot{},
		&models.Server{},
		&models.UserBot{},
		&models.BotPlanet{},
		&models.Test{},
	)

	return db
}
