package service

import (
	"github.com/whyslove/avito-test/core/models"
	"github.com/whyslove/avito-test/core/repository"
)

type User interface {
	GetBalance(int) (float64, error) // int is id
	AddMoney(models.UserOperation) error
	WriteOffMoney(models.UserOperation) error
	TransferOperation(models.UserTransferOperation) error
}

type Service struct {
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos),
	}
}
