package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Operation struct {
	Id             int       `json:"id"`
	MainSubjectId  int       `json:"main_subject_id"`
	OtherSubjectId NullInt   `json:"other_subject_id"`
	ExecutedAt     time.Time `json:"executed_at"`
	Money          float64   `json:"amount_of_money"`
}

type NullInt struct {
	sql.NullInt32
}

func (ni *NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int32)
}

func (ni NullInt) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int32)
	ni.Valid = (err == nil)
	return err
}
