package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict raised")
)

type Models struct {
	Movies EpisodeModel
	Tokens TokenModel

	Users UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: EpisodeModel{DB: db},
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
