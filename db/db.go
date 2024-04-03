package db

import (
	"log"
	"oude-bar-picker-v2/model"

	libsql "github.com/renxzen/gorm-libsql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(libsql.Open("libsql://oude-picker-mathiasare.turso.io?authToken=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhIjoicnciLCJpYXQiOjE3MTA3OTgzOTEsImlkIjoiZTVmZGI3MzctYzJlYy0xMWVlLWIwNzEtNTIwNTNjNGIyM2E2In0.nLQd23Cua7IC2tfUyWFg3l1APm5BTZMzguF_tAJGupRsB4tlKuvw5_n2aiCPcP2kGQ51hhW3nG0tFBo3Ag4TAQ"), &gorm.Config{})

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
