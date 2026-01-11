package db

import (
	"log"

	"gorm.io/driver/postgres" // or mysql, sqlite, etc.
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&User{},
		&File{},
		&UniqueFace{},
		&Face{},
		&Job{},
	)

	if err != nil {
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}
