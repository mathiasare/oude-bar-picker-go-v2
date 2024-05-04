package service

import (
	"errors"
	"oude-bar-picker-v2/model"
	"oude-bar-picker-v2/repository"
	"sort"

	"gorm.io/gorm"
)

type ParticipantService struct {
	Repo repository.ParticipantRepository
}

func (service ParticipantService) getRepository() repository.ParticipantRepository {
	return service.Repo
}

func (service ParticipantService) GetAll() ([]model.Participant, error) {
	repo := service.getRepository()
	return repo.FindAll()
}

func (service ParticipantService) GetById(id uint) (model.Participant, error) {
	repo := service.getRepository()
	return repo.FindById(id)
}

func (service ParticipantService) Create(p model.Participant) (model.Participant, error) {
	repo := service.getRepository()
	return repo.Save(p)
}

func (service ParticipantService) Delete(id uint) (uint, error) {
	repo := service.getRepository()
	return repo.DeleteById(id)
}

func (service ParticipantService) ExistsNameAndVote(name string, voteCode string) (bool, error) {
	repo := service.getRepository()
	_, err := repo.FindByNameAndVote(name, voteCode)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (service ParticipantService) GetByNameAndVote(name string, voteCode string) (model.Participant, error) {
	repo := service.getRepository()
	p, err := repo.FindByNameAndVote(name, voteCode)

	if err != nil {
		return model.Participant{}, err
	}
	return p, nil
}

func (service ParticipantService) VoteForBar(p *model.Participant, barId uint) ([]model.Participant, error) {
	repo := service.getRepository()
	res, err := repo.UpdateBarId(p, barId)

	if err != nil {
		return nil, err
	}

	return repo.FindAllForVoteWhereBarNotNull(res.VoteId)
}

func (service ParticipantService) GetAllVotedParticipants(voteCode string) ([]model.Participant, error) {
	repo := service.getRepository()
	return repo.FindAllForVoteWhereBarNotNull(voteCode)
}

func (service ParticipantService) GetAllParticipantsForVote(voteCode string) ([]model.Participant, error) {
	repo := service.getRepository()
	return repo.FindAllForVote(voteCode)
}

func (service ParticipantService) GetVoteStats(votedPs []model.Participant) model.VoteStatsDTO {
	statsMap := make(map[uint]*model.VoteStatsRow)
	rows := []model.VoteStatsRow{}

	for _, p := range votedPs {
		if p.BarId == nil {
			continue
		}

		row, exists := statsMap[*p.BarId]

		if !exists {
			row = &model.VoteStatsRow{BarId: *p.BarId, BarName: p.Bar.Name, VoteCount: 0}
			statsMap[*p.BarId] = row
		}

		row.VoteCount = row.VoteCount + 1
	}

	for _, row := range statsMap {
		rows = append(rows, *row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].VoteCount > rows[j].VoteCount
	})
	return rows
}
