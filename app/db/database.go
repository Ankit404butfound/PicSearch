package db

import (
	"log"
	"os"

	"PicSearch/app/db/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres" // or mysql, sqlite, etc.
	"gorm.io/gorm"
)

var err = godotenv.Load(".env")

var DB, _ = gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
	DisableForeignKeyConstraintWhenMigrating: false,
})

func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.UniqueFace{},
		&models.Face{},
		&models.Job{},
		&models.Devices{},
	)

	if err != nil {
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}
