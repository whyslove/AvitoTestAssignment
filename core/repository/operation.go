package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/whyslove/avito-test/core/models"
)

type OperationPostgres struct {
	db *sqlx.DB
}

func NewOperationPostgres(db *sqlx.DB) *OperationPostgres {
	return &OperationPostgres{db: db}
}

func (op *OperationPostgres) InsertOperation(Tx *sql.Tx, operation models.Operation) error {
	query := fmt.Sprintf("INSERT INTO %s (main_subject_id, other_subject_id, amount_of_money, executed_at) VALUES ($1, $2, $3, $4);", operationTable)
	_, err := Tx.Exec(query, operation.MainSubjectId, operation.OtherSubjectId, operation.Money, operation.ExecutedAt)
	if err != nil {
		return err
	}
	return nil
}
func (op *OperationPostgres) GetUserOperations(Tx *sql.Tx, userId, lower_bound, upper_bound int) ([]models.Operation, error) {
	var operations []models.Operation
	query := fmt.Sprintf("SELECT * FROM %s WHERE main_subject_id = $1 LIMIT %d OFFSET %d",
		operationTable, upper_bound-lower_bound, lower_bound)
	rows, err := Tx.Query(query, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var o models.Operation
		if err := rows.Scan(&o.Id, &o.MainSubjectId, &o.OtherSubjectId, &o.Money, &o.ExecutedAt); err != nil {
			return nil, err
		}
		logrus.Debug(o)
		operations = append(operations, o)
	}
	return operations, nil

}
func (op *OperationPostgres) GetUserOperationsSorted(Tx *sql.Tx, userId, lowerBound, upperBound int, sortType string) ([]models.Operation, error) {
	var operations []models.Operation
	var query string
	switch sortType {
	case "summ":
		query = fmt.Sprintf(`SELECT * FROM %s WHERE main_subject_id = $1
		ORDER BY amount_of_money DESC
		LIMIT %d OFFSET %d`,
			operationTable, upperBound-lowerBound, lowerBound)
	case "date":
		query = fmt.Sprintf(`SELECT * FROM %s WHERE main_subject_id = $1
		ORDER BY executed_at DESC
		LIMIT %d OFFSET %d`,
			operationTable, upperBound-lowerBound, lowerBound)
	default:
		return nil, fmt.Errorf("do not recognise sort type %s", sortType)
	}
	rows, err := Tx.Query(query, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var o models.Operation
		if err := rows.Scan(&o.Id, &o.MainSubjectId, &o.OtherSubjectId, &o.Money, &o.ExecutedAt); err != nil {
			return nil, err
		}
		logrus.Debug(o)
		operations = append(operations, o)
	}
	return operations, nil
}
