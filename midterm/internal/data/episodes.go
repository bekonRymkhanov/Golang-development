package data

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"series.bekarysrymkhanov.net/internal/validator"
	"time"
)

type Episodes struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"-"`
	Title      string    `json:"title"`
	Year       int32     `json:"year,omitempty"`
	Runtime    Runtime   `json:"runtime,omitempty"`
	Characters []string  `json:"characters,omitempty"`
	Version    int32     `json:"version"`
}

type EpisodeModel struct {
	DB *sql.DB
}

func (e EpisodeModel) Insert(episode *Episode) error {
	query := `INSERT INTO episodes (title, year, runtime, characters) 
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at, version`

	args := []interface{}{episode.Title, episode.Year, episode.Runtime, pq.Array(episode.Characters)}

	return e.DB.QueryRow(query, args...).Scan(&episode.ID, &episode.CreatedAt, &episode.Version)
}

func (e EpisodeModel) Get(id int64) (*Episode, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, year, runtime, characters, version
				FROM episodes
				WHERE id = $1`
	var episode Episode
	err := e.DB.QueryRow(query, id).Scan(
		&episode.ID,
		&episode.CreatedAt,
		&episode.Title,
		&episode.Year,
		&episode.Runtime,
		pq.Array(&episode.Characters),
		&episode.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &episode, nil

}
func (e EpisodeModel) Update(episod *Episode) error {
	query := `UPDATE episodes
				SET title = $1, year = $2, runtime = $3, characters = $4, version = version + 1
				WHERE id = $5
				RETURNING version`

	args := []interface{}{
		episod.Title,
		episod.Year,
		episod.Runtime,
		pq.Array(episod.Characters),
		episod.ID,
	}
	err := e.DB.QueryRow(query, args...).Scan(&episod.Version)
	return err
}
func (e EpisodeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM episodes
				WHERE id=$1`

	result, err := e.DB.Exec(query, id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return err
}

func ValidateMovie(v *validator.Validator, episode *Episode) {
	v.Check(episode.Title != "", "title", "must be provided")
	v.Check(len(episode.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(episode.Year != 0, "year", "must be provided")
	v.Check(episode.Year >= 1888, "year", "must be greater than 1888")
	v.Check(episode.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(episode.Runtime != 0, "runtime", "must be provided")
	v.Check(episode.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(episode.Characters != nil, "cahracters", "must be provided")
	v.Check(len(episode.Characters) >= 1, "characters", "must contain at least 1 characters")
	v.Check(len(episode.Characters) <= 20, "characters", "must not contain more than 20 characters")
	v.Check(validator.Unique(episode.Characters), "genres", "must not contain duplicate values")
}
