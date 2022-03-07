package service

import (
	"fmt"
	"time"

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

	var operation models.Operation
	operation.MainSubjectId = uto.SenderId
	operation.OtherSubjectId.Int32 = int32(uto.ReceiverId)
	operation.OtherSubjectId.Valid = true
	operation.Money = uto.AmountOfMoney
	operation.ExecutedAt = time.Now().UTC()

	err3 := s.repo.InsertOperation(Tx, operation)

	//Next reverse some fields to store operation to other user
	operation.MainSubjectId = uto.ReceiverId
	operation.OtherSubjectId.Int32 = int32(uto.SenderId)
	operation.Money = -operation.Money

	err4 := s.repo.InsertOperation(Tx, operation)

	if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
		s.repo.CommitTransaction(Tx)
	} else {
		return fmt.Errorf("errors in updating blanace, err_send: %s, err_rec: %s, err_from_add_info %s and %s",
			err1.Error(), err2.Error(), err3.Error(), err4.Error())
	}
	return nil
}

func (s *UserService) AddMoney(uo models.UserOperation) error {
	Tx, err := s.repo.StartTransaction()
	if err != nil {
		return err
	}
	defer s.repo.RollbackTransaction(Tx)
	balance, err := s.repo.GetBalanceTx(Tx, uo.UserId)
	if err != nil {
		if err == models.NoRecordInDb {
			err = s.repo.InsertUserTx(Tx, models.User{Id: uo.UserId, Balance: uo.AmountOfMoney})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	err1 := s.repo.UpdateUserTx(Tx,
		models.User{Id: uo.UserId, Balance: balance + uo.AmountOfMoney})
	err2 := s.repo.InsertOperation(Tx, models.Operation{MainSubjectId: uo.UserId, ExecutedAt: time.Now().UTC(), Money: uo.AmountOfMoney})
	if err1 == nil && err2 == nil {
		s.repo.CommitTransaction(Tx)
	} else {
		return fmt.Errorf("errors in updating blanace, err_send: %s,err_from_add_info%s",
			err1.Error(), err2.Error())
	}
	return nil
}

func (s *UserService) WriteOffMoney(uo models.UserOperation) error {
	Tx, err := s.repo.StartTransaction()
	if err != nil {
		return err
	}
	defer s.repo.RollbackTransaction(Tx)
	balance, err := s.repo.GetBalanceTx(Tx, uo.UserId)
	if err != nil {
		return err
	}
	if balance+uo.AmountOfMoney < 0 {
		return fmt.Errorf("not enough money for this operation")
	} //amount of money always < 0
	err1 := s.repo.UpdateUserTx(Tx, models.User{Id: uo.UserId, Balance: balance + uo.AmountOfMoney})
	err2 := s.repo.InsertOperation(Tx, models.Operation{MainSubjectId: uo.UserId, ExecutedAt: time.Now().UTC(), Money: uo.AmountOfMoney})
	if err1 == nil && err2 == nil {
		s.repo.CommitTransaction(Tx)
	} else {
		return fmt.Errorf("errors in updating blanace, err_send: %s,err_from_add_info%s",
			err1.Error(), err2.Error())
	}
	return nil
}
