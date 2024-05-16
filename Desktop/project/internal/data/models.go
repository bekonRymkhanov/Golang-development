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
	Movies      EpisodeModel
	Characters  CharacterModel
	Tokens      TokenModel
	Permissions PermissionModel
	Users       UserModel
	LikeComment LikeCommentModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:     EpisodeModel{DB: db},
		Characters: CharacterModel{DB: db},
		LikeComment: LikeCommentModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
