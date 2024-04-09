package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"series.bekarysrymkhanov.net/internal/validator"
	"time"
)

type EpisodeModel struct {
	DB *sql.DB
}

func (e EpisodeModel) Insert(episode *Episode) error {
	query := `INSERT INTO episodes (title, year, runtime, characters) 
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at, version`

	args := []interface{}{episode.Title, episode.Year, episode.Runtime, pq.Array(episode.Characters)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&episode.ID, &episode.CreatedAt, &episode.Version)
}

func (e EpisodeModel) Get(id int64) (*Episode, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, year, runtime, characters, version
				FROM episodes
				WHERE id = $1`
	var episode Episode

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, id).Scan(

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
				WHERE id = $5 and version = $6
				RETURNING version`

	args := []interface{}{
		episod.Title,
		episod.Year,
		episod.Runtime,
		pq.Array(episod.Characters),
		episod.ID,
		episod.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&episod.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return err
}
func (e EpisodeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM episodes
				WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := e.DB.ExecContext(ctx, query, id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return err
}

func (e EpisodeModel) GetAll(title string, characters []string, filters Filters) ([]*Episode, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, characters, version
		FROM episodes
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (characters @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query, title, pq.Array(characters), filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	episodes := []*Episode{}

	for rows.Next() {
		var episode Episode

		err := rows.Scan(
			&totalRecords,
			&episode.ID,
			&episode.CreatedAt,
			&episode.Title,
			&episode.Year,
			&episode.Runtime,
			pq.Array(&episode.Characters),
			&episode.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		episodes = append(episodes, &episode)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return episodes, metadata, nil
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
