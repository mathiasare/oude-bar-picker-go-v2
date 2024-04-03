package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gorm.io/gorm"

	"oude-bar-picker-v2/model"
)

func Seed(db *gorm.DB) {
	jsonFile, err := os.Open("./resource/bars.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var bars []model.Bar

	json.Unmarshal(byteValue, &bars)
	db.Create(&bars)
}
