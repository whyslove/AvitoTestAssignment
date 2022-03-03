package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/whyslove/avito-test/core/models"
)

type User interface {
	GetBalance(id int) (float64, error)
	GetBalanceTx(Tx *sql.Tx, id int) (float64, error)
	InsertUser(us models.User) error
	UpdateUser(us models.User) error
	InsertUserTx(Tx *sql.Tx, us models.User) error
	UpdateUserTx(Tx *sql.Tx, us models.User) error
	StartTransaction() (*sql.Tx, error)
	CommitTransaction(Tx *sql.Tx) error
	RollbackTransaction(Tx *sql.Tx) error
}

type Repository struct {
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserPostgres(db),
	}
}
