package main

import (
	"PicSearch/app/db"
	"log"
)

func main() {
	// Database connection string
	dsn := "postgresql://rpie:rpie@100.115.44.83:5432/picsearch"

	// Connect to database
	err := db.ConnectDatabase(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	err = db.Migrate()
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration completed!")
}
