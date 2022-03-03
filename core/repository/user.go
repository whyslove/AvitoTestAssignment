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

func (u *UserPostgres) InsertUser(us models.User) error {
	query := fmt.Sprintf("INSERT INTO %s (id, balance) VALUES ($1, $2);", userTable)
	_, err := u.db.Exec(query, us.Id, us.Balance)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserPostgres) UpdateUser(us models.User) error {
	query := fmt.Sprintf("UPDATE %s SET balance=$1 WHERE id=$2", userTable)
	_, err := u.db.Exec(query, us.Balance, us.Id)
	if err != nil {
		return err
	}
	return nil
}
