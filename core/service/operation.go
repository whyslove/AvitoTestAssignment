package service

import (
	"github.com/whyslove/avito-test/core/models"
	"github.com/whyslove/avito-test/core/repository"
)

const pageLen = 10

type OperationService struct {
	repo *repository.Repository
}

func NewOperationService(repo *repository.Repository) *OperationService {
	return &OperationService{repo: repo}
}

func (os *OperationService) GetOperations(user_id int, page int) ([]models.Operation, error) {
	Tx, err := os.repo.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer os.repo.RollbackTransaction(Tx)
	lowerBound := (page - 1) * pageLen
	upperBound := page * pageLen
	operations, err := os.repo.GetUserOperations(Tx, user_id, lowerBound, upperBound)
	if err != nil {
		return nil, err
	}
	return operations, nil
}
func (os *OperationService) GetSortedOperations(user_id int, page int, sortType string) ([]models.Operation, error) {
	Tx, err := os.repo.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer os.repo.RollbackTransaction(Tx)
	lowerBound := (page - 1) * pageLen
	upperBound := page * pageLen
	operations, err := os.repo.GetUserOperationsSorted(Tx, user_id, lowerBound, upperBound, sortType)
	if err != nil {
		return nil, err
	}
	return operations, nil
}
