package service

import (
	"encoding/json"
	"io"
	"oude-bar-picker-v2/model"
	"oude-bar-picker-v2/repository"
)

type BarService struct {
	Repo repository.BarRepository
}

func (service BarService) getRepository() repository.BarRepository {
	return service.Repo
}

func (service BarService) GetAll() ([]model.Bar, error) {
	repo := service.getRepository()
	return repo.FindAll()
}

func (service BarService) GetById(id uint) (model.Bar, error) {
	repo := service.getRepository()
	return repo.FindById(id)
}

func (service BarService) Create(body io.ReadCloser) (model.Bar, error) {
	var bar model.Bar
	json.NewDecoder(body).Decode(&bar)
	repo := service.getRepository()
	return repo.Save(bar)
}

func (service BarService) Delete(id string) (string, error) {
	repo := service.getRepository()
	return repo.DeleteById(id)
}
