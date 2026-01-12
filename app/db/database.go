package db

import (
	"log"

	"PicSearch/app/db/models"

	"gorm.io/driver/postgres" // or mysql, sqlite, etc.
	"gorm.io/gorm"
)

var DB, err = gorm.Open(postgres.Open("postgresql://rpie:rpie@100.115.44.83:5432/picsearch"), &gorm.Config{
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
