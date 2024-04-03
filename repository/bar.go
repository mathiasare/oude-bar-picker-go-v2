package repository

import (
	"oude-bar-picker-v2/model"

	"gorm.io/gorm"
)

type BarRepository struct {
	Db *gorm.DB
}

func (repo BarRepository) getDB() *gorm.DB {
	return repo.Db
}

func (repo BarRepository) FindAll() ([]model.Bar, error) {
	db := repo.getDB()
	var bars []model.Bar
	result := db.Find(&bars)

	err := result.Error
	if err != nil {
		return []model.Bar{}, err
	}

	return bars, nil
}

func (repo BarRepository) FindById(id uint) (model.Bar, error) {
	db := repo.getDB()
	var bar model.Bar
	result := db.First(&bar, id)

	err := result.Error
	if err != nil {
		return model.Bar{}, err
	}

	return bar, nil
}

func (repo BarRepository) Save(bar model.Bar) (model.Bar, error) {
	db := repo.getDB()
	result := db.Create(&bar)

	err := result.Error
	if err != nil {
		return model.Bar{}, err
	}

	return bar, nil
}

func (repo BarRepository) DeleteById(id string) (string, error) {
	db := repo.getDB()
	result := db.Delete(&model.Bar{}, id)

	err := result.Error
	if err != nil {
		return "-1", err
	}

	return id, nil
}
