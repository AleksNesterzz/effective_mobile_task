package storage

import (
	"fmt"
	"future_today/internal/cerrors"
	"future_today/internal/config"
	"future_today/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb(cfg *config.Config) (*gorm.DB, error) {
	dbConnString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DbHost, cfg.DbUser, cfg.DbPass, cfg.DbName, cfg.DbPort,
	)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		return nil, cerrors.ErrDbConnect
	}

	err = db.AutoMigrate(&models.Person{})
	if err != nil {
		return nil, cerrors.ErrMigration
	}


	return db, nil
}
