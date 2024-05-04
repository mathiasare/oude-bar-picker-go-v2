package repository

import (
	"oude-bar-picker-v2/model"

	"gorm.io/gorm"
)

type ParticipantRepository struct {
	Db *gorm.DB
}

func (repo ParticipantRepository) getDB() *gorm.DB {
	return repo.Db
}

func (repo ParticipantRepository) FindAll() ([]model.Participant, error) {
	db := repo.getDB()
	var participants []model.Participant
	result := db.Find(&participants)

	err := result.Error
	if err != nil {
		return []model.Participant{}, err
	}

	return participants, nil
}

func (repo ParticipantRepository) FindAllForVote(voteCode string) ([]model.Participant, error) {
	db := repo.getDB()
	var participants []model.Participant
	result := db.Preload("Bar").Where("vote_id = ?", voteCode).Find(&participants)

	err := result.Error
	if err != nil {
		return nil, err
	}

	return participants, nil
}

func (repo ParticipantRepository) FindAllForVoteWhereBarNotNull(voteCode string) ([]model.Participant, error) {
	db := repo.getDB()
	var participants []model.Participant
	result := db.Preload("Bar").Where("bar_id IS NOT NULL AND vote_id = ?", voteCode).Find(&participants)

	err := result.Error
	if err != nil {
		return nil, err
	}

	return participants, nil
}

func (repo ParticipantRepository) FindById(id uint) (model.Participant, error) {
	db := repo.getDB()
	var participant model.Participant
	result := db.First(&participant, id)

	err := result.Error
	if err != nil {
		return model.Participant{}, err
	}

	return participant, nil
}

func (repo ParticipantRepository) FindByNameAndVote(name string, voteCode string) (model.Participant, error) {
	db := repo.getDB()

	var participant model.Participant
	result := db.Where("name = ? AND vote_id = ?", name, voteCode).First(&participant)

	err := result.Error
	if err != nil {
		return model.Participant{}, err
	}

	return participant, nil
}

func (repo ParticipantRepository) Save(participant model.Participant) (model.Participant, error) {
	db := repo.getDB()
	result := db.Create(&participant)

	err := result.Error
	if err != nil {
		return model.Participant{}, err
	}

	return participant, nil
}

func (repo ParticipantRepository) DeleteById(id uint) (uint, error) {
	db := repo.getDB()
	result := db.Delete(&model.Participant{}, id)

	err := result.Error
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo ParticipantRepository) UpdateBarId(p *model.Participant, barId uint) (model.Participant, error) {
	tx := repo.getDB().Model(&p).Update("bar_id", barId)
	err := tx.Error
	if err != nil {
		return model.Participant{}, err
	}
	return repo.FindById(p.ID)
}
