package Complite

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func init() {
	initDB()
}

func initDB() {
	dsn := "host=localhost user=postgres dbname=postgres port=5432 sslmode=disable password=mysecretpassword"
	//Angular-postgres
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&UserRegister{}, &ChatMessage{})
	log.Println("Database successfully connected and migrated is Angular!")
}
