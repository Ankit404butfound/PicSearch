package db

import (
	"log"
	"os"

	"PicSearch/app/db/models"

	"gorm.io/driver/postgres" // or mysql, sqlite, etc.
	"gorm.io/gorm"
)

var DB, err = gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
	DisableForeignKeyConstraintWhenMigrating: false,
})

func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.UniqueFace{},
		&models.Face{},
		&models.Job{},
	)

	if err != nil {
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}
