package service

import (
	"evo-test/internal/models"
	"evo-test/internal/repository"
)

type Service interface {
	GetTransaction(params models.SearchParams) ([]models.Transaction, error)
	LoadData(transactions []models.Transaction) error
}

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTransaction(params models.SearchParams) ([]models.Transaction, error) {
	return s.GetTransaction(params)
}

func (s *service) LoadData(transactions []models.Transaction) error {
	return s.repo.InsertData(transactions)
}
