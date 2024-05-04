package repository

import (
	"oude-bar-picker-v2/model"

	"gorm.io/gorm"
)

type VoteRepository struct {
	Db *gorm.DB
}

func (repo VoteRepository) getDB() *gorm.DB {
	return repo.Db
}

func (repo VoteRepository) FindAll() ([]model.Vote, error) {
	db := repo.getDB()
	var votes []model.Vote
	result := db.Find(&votes)

	err := result.Error
	if err != nil {
		return []model.Vote{}, err
	}

	return votes, nil
}

func (repo VoteRepository) FindById(id string) (model.Vote, error) {
	db := repo.getDB()
	var vote model.Vote
	result := db.Preload("Participants").First(&vote, "id = ?", id)

	err := result.Error
	if err != nil {
		return model.Vote{}, err
	}

	return vote, nil
}

func (repo VoteRepository) FindDeletedById(id string) (model.Vote, error) {
	db := repo.getDB()
	var vote model.Vote
	result := db.Unscoped().First(&vote, "id = ?", id)

	err := result.Error
	if err != nil {
		return model.Vote{}, err
	}

	return vote, nil
}

func (repo VoteRepository) Save(vote model.Vote) (model.Vote, error) {
	db := repo.getDB()
	result := db.Create(&vote)

	err := result.Error
	if err != nil {
		return model.Vote{}, err
	}

	return vote, nil
}

func (repo VoteRepository) DeleteById(id string) (string, error) {
	db := repo.getDB()
	result := db.Delete(&model.Vote{}, "id = ?", id)

	err := result.Error
	if err != nil {
		return "-1", err
	}

	return id, nil
}

func (repo VoteRepository) AddWinnerToVote(id string, winnerId uint) error {
	tx := repo.getDB().Model(&model.Vote{}).Where("id = ?", id).Update("winner_id", winnerId)
	err := tx.Error
	if err != nil {
		return err
	}
	return nil
}
