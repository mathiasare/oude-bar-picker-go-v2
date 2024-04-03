package service

import (
	"log"
	"oude-bar-picker-v2/model"
	"oude-bar-picker-v2/repository"

	"github.com/google/uuid"
)

type VoteService struct {
	PService ParticipantService
	Repo     repository.VoteRepository
}

func (service VoteService) getRepository() repository.VoteRepository {
	return service.Repo
}

func (service VoteService) GetAll() ([]model.Vote, error) {
	repo := service.getRepository()
	return repo.FindAll()
}

func (service VoteService) GetById(voteCode string) (model.Vote, error) {
	repo := service.getRepository()
	return repo.FindById(voteCode)
}

func (service VoteService) Create(vote model.Vote) (model.Vote, error) {
	repo := service.getRepository()
	voteCode := generateVoteCode()
	vote.ID = voteCode
	return repo.Save(vote)
}

func (service VoteService) Delete(voteCode string) (string, error) {
	repo := service.getRepository()
	return repo.DeleteById(voteCode)
}

func (service VoteService) AddUserToVote(voteCode string, name string) error {
	pService := service.PService

	existsUser, err := pService.ExistsNameAndVote(name, voteCode)
	if err != nil {
		return err
	}

	if existsUser {
		log.Printf("Participant %s already exists for vote %s", name, voteCode)
		return nil
	}

	_, err = pService.Create(model.Participant{Name: name, VoteId: voteCode})
	if err != nil {
		return err
	}

	return nil
}

func (service VoteService) EndVote(voteCode string, name string) (model.VoteStatsDTO, error) {
	pService := service.PService
	repo := service.getRepository()
	votedPs, err := pService.GetAllVotedParticipants(voteCode)
	if err != nil {
		return nil, err
	}

	stats := pService.GetVoteStats(votedPs)
	if len(stats) == 0 {
		_, err := repo.DeleteById(voteCode)
		if err != nil {
			return nil, err
		}
		return stats, nil
	}

	repo.AddWinnerToVote(voteCode, stats[0].BarId)

	return stats, nil
}

func generateVoteCode() string {
	return uuid.New().String()[:6]
}
