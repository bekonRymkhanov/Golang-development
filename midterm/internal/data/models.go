package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Movies EpisodeModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: EpisodeModel{DB: db},
	}
}
