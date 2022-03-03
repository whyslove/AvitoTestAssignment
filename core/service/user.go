package service

import (
	"fmt"

	"github.com/whyslove/avito-test/core/models"
	"github.com/whyslove/avito-test/core/repository"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetBalance(id int) (float64, error) {
	return s.repo.GetBalance(id)
}
func (s *UserService) TransferOperation(uto models.UserTransferOperation) error {
	Tx, err := s.repo.StartTransaction()
	if err != nil {
		return err
	}
	defer s.repo.RollbackTransaction(Tx)
	sender_balance, err := s.repo.GetBalanceTx(Tx, uto.SenderId)
	if err != nil {
		return err
	}
	receiver_balance, err := s.repo.GetBalanceTx(Tx, uto.ReceiverId)
	if err != nil {
		if err != models.NoRecordInDb {
			return err
		} else {
			err = s.repo.InsertUserTx(Tx, models.User{Id: uto.ReceiverId, Balance: 0})
			if err != nil {
				return err
			}
			receiver_balance = 0
		}

	}

	if sender_balance-uto.AmountOfMoney < 0 {
		return fmt.Errorf("id=%d has not enough money ", uto.SenderId)
	}
	err1 := s.repo.UpdateUserTx(Tx, models.User{Id: uto.SenderId, Balance: sender_balance - uto.AmountOfMoney})
	err2 := s.repo.UpdateUserTx(Tx, models.User{Id: uto.ReceiverId, Balance: receiver_balance + uto.AmountOfMoney})
	if err1 == nil && err2 == nil {
		s.repo.CommitTransaction(Tx)
	} else {
		return fmt.Errorf("errors in updating blanace, err_send: %s, err_rec: %s", err1.Error(), err2.Error())
	}
	return nil
}

func (s *UserService) AddMoney(uo models.UserOperation) error {
	balance, err := s.repo.GetBalance(uo.UserId)
	if err != nil {
		if err == models.NoRecordInDb {
			err = s.repo.InsertUser(models.User{Id: uo.UserId, Balance: uo.AmountOfMoney})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return s.repo.UpdateUser(
			models.User{Id: uo.UserId, Balance: balance + uo.AmountOfMoney})
	}
	return nil
}

func (s *UserService) WriteOffMoney(uo models.UserOperation) error {
	balance, err := s.repo.GetBalance(uo.UserId)
	if err != nil {
		return err
	}
	if balance+uo.AmountOfMoney < 0 {
		return fmt.Errorf("not enough money for this operation")
	} //amount of money always < 0
	return s.repo.UpdateUser(models.User{Id: uo.UserId, Balance: balance + uo.AmountOfMoney})

}
