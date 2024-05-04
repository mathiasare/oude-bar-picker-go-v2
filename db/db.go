package db

import (
	"log"
	"os"
	"oude-bar-picker-v2/model"

	libsql "github.com/renxzen/gorm-libsql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dbUrl := os.Getenv("DB_URL")
	db, err := gorm.Open(libsql.Open(dbUrl), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&model.Bar{}, &model.Vote{}, &model.Participant{})
	if err != nil {
		log.Panic(err, "Failed to auto migrate")
	}

	log.Println("Successfully connected to database!")
	return db
}
