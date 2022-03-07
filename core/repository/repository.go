package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/whyslove/avito-test/core/models"
)

type User interface {
	GetBalance(id int) (float64, error)
	GetBalanceTx(Tx *sql.Tx, id int) (float64, error)
	InsertUserTx(Tx *sql.Tx, us models.User) error
	UpdateUserTx(Tx *sql.Tx, us models.User) error
	StartTransaction() (*sql.Tx, error)
	CommitTransaction(Tx *sql.Tx) error
	RollbackTransaction(Tx *sql.Tx) error
}
type Operation interface {
	GetUserOperations(Tx *sql.Tx, userId, lower_bound, upper_bound int) ([]models.Operation, error)
	GetUserOperationsSorted(Tx *sql.Tx, userId, lowerBound, upperBound int, sortType string) ([]models.Operation, error)
	InsertOperation(Tx *sql.Tx, operation models.Operation) error
}

type Repository struct {
	User
	Operation
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:      NewUserPostgres(db),
		Operation: NewOperationPostgres(db),
	}
}
