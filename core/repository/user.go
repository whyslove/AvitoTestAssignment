package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/whyslove/avito-test/core/models"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}
func (u *UserPostgres) StartTransaction() (*sql.Tx, error) {
	return u.db.Begin()
}

func (u *UserPostgres) GetBalance(id int) (float64, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", userTable)
	if err := u.db.Get(&user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return 0, models.NoRecordInDb
		} else {
			logrus.Debug("other error")
		}
	}
	logrus.Debug(user)
	return user.Balance, nil
}

func (u *UserPostgres) GetBalanceTx(Tx *sql.Tx, id int) (float64, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", userTable)
	row := Tx.QueryRow(query, id)
	if err := row.Scan(&user.Id, &user.Balance); err != nil {
		if err == sql.ErrNoRows {
			return 0, models.NoRecordInDb
		} else {
			logrus.Debug("other error in db")
			return 0, err
		}
	}
	logrus.Debug(user)
	return user.Balance, nil
}

func (u *UserPostgres) InsertUserTx(Tx *sql.Tx, us models.User) error {
	query := fmt.Sprintf("INSERT INTO %s (id, balance) VALUES ($1, $2);", userTable)
	_, err := Tx.Exec(query, us.Id, us.Balance)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserPostgres) UpdateUserTx(Tx *sql.Tx, us models.User) error {
	query := fmt.Sprintf("UPDATE %s SET balance=$1 WHERE id=$2", userTable)
	_, err := Tx.Exec(query, us.Balance, us.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserPostgres) CommitTransaction(Tx *sql.Tx) error {
	if err := Tx.Commit(); err != nil {
		return fmt.Errorf("error in commit transaction")
	}
	return nil
}

func (u *UserPostgres) RollbackTransaction(Tx *sql.Tx) error {
	return Tx.Rollback()

}
