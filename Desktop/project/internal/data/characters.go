package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"series.bekarysrymkhanov.net/internal/validator"
	"time"
)

type CharacterModel struct {
	DB *sql.DB
}

func (e CharacterModel) Insert(character *Character) error {
	query := `INSERT INTO characters (episodes_id, name, age) 
				VALUES ($1, $2, $3)
				RETURNING id, version`

	args := []interface{}{character.EpisodesID, character.Name, character.Age}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&character.ID, &character.Version)
}
func (e CharacterModel) Get(id int64) (*Character, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, episodes_id, name, age, version
				FROM characters
				WHERE id = $1`
	var character Character

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, id).Scan(
		&character.ID,
		&character.EpisodesID,
		&character.Name,
		&character.Age,
		&character.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &character, nil
}

func (e CharacterModel) Update(character *Character) error {
	query := `UPDATE characters
				SET episodes_id = $1, name = $2, age = $3, version = version + 1
				WHERE id = $4 AND version = $5
				RETURNING version`

	args := []interface{}{
		character.EpisodesID,
		character.Name,
		character.Age,
		character.ID,
		character.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&character.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (e CharacterModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM characters
				WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := e.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (e CharacterModel) GetAll(name string, filters Filters) ([]*Character, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, episodes_id, name, age, version
		FROM characters
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query, name, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	characters := []*Character{}

	for rows.Next() {
		var character Character
		err := rows.Scan(
			&totalRecords,
			&character.ID,
			&character.EpisodesID,
			&character.Name,
			&character.Age,
			&character.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		characters = append(characters, &character)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return characters, metadata, nil
}
func (e CharacterModel) GetByEpisodeID(episodeID int64) ([]*Character, error) {
	query := `
		SELECT id, episodes_id, name, age, version
		FROM characters
		WHERE episodes_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query, episodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*Character

	for rows.Next() {
		var character Character
		err := rows.Scan(
			&character.ID,
			&character.EpisodesID,
			&character.Name,
			&character.Age,
			&character.Version,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, &character)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil
}

func ValidateCharacter(v *validator.Validator, character *Character) {
	v.Check(character.Name != "", "name", "must be provided")
	v.Check(len(character.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(character.Age != 0, "year", "must be provided")
	v.Check(character.Age <= 1000, "year", "must be less than 1000")
}
