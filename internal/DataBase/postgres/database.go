package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB

func init() {
	initDB()
}

func initDB() {
	dsn := "host=localhost user=postgres dbname=postgres port=5432 sslmode=disable password=mysecretpassword"
	//Angular-postgres
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	Db.AutoMigrate(&UserRegister{}, &ChatMessage{})
	log.Println("Database successfully connected and migrated is Angular!")
}
