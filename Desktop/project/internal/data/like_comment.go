package data

import (
	"context"
	"database/sql"
	//"errors"
	"fmt"
	"log"
	"time"
	//"github.com/lib/pq"
	"series.bekarysrymkhanov.net/internal/validator"
)

type LikeComment struct {
	LikeID      int       `json:"id"`
	UserID      int       `json:"user_id"`
	EpisodeID   int       `json:"episode_id"`
	CommentText string    `json:"comment_text"`
	LikeCount   int       `json:"like_count"`
	CreatedAt   time.Time `json:"created_at"`
	Version     int       `json:"-"`
}
type LikeCommentModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (lcm *LikeCommentModel) Insert(likeComment *LikeComment) error {
	query := `INSERT INTO like_comment (user_id, episode_id, comment_text) 
				VALUES ($1, $2, $3)
				RETURNING id,created_at`

	args := []interface{}{likeComment.UserID, likeComment.EpisodeID, likeComment.CommentText}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return lcm.DB.QueryRowContext(ctx, query, args...).Scan(&likeComment.LikeID, &likeComment.CreatedAt)
}

func (lcm *LikeCommentModel) Delete(likeID int64) error {
	if likeID < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM like_comment
				WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := lcm.DB.ExecContext(ctx, query, likeID)
	return err
}

func (lcm *LikeCommentModel) Update(likeComment *LikeComment) error {
	query := `UPDATE like_comment
				SET comment_text = $1, like_count = $2, version = version + 1
				WHERE id = $3 AND version = $4
				RETURNING version`

	args := []interface{}{
		likeComment.CommentText,
		likeComment.LikeCount,
		likeComment.LikeID,
		likeComment.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//err := lcm.DB.QueryRowContext(ctx, query, args...).Scan(&likeComment.Version)
	return lcm.DB.QueryRowContext(ctx, query, args...).Scan(&likeComment.Version)
}

func (lcm *LikeCommentModel) Get(id int64) (*LikeComment, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, user_id, episode_id, comment_text, like_count, created_at, version
				FROM like_comment
				WHERE id = $1`

	var likeComment LikeComment

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := lcm.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&likeComment.LikeID,
		&likeComment.UserID,
		&likeComment.EpisodeID,
		&likeComment.CommentText,
		&likeComment.LikeCount,
		&likeComment.CreatedAt,
		&likeComment.Version,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive like_comment with id: %v, %w", id, err)
	}
	return &likeComment, nil
}

func (e LikeCommentModel) GetAll(commentText string, filters Filters) ([]*LikeComment, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, user_id, episode_id, comment_text, like_count, created_at, version
		FROM like_comment
		WHERE (to_tsvector('simple',comment_text) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query, commentText, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	likes := []*LikeComment{}

	for rows.Next() {
		var like LikeComment

		err := rows.Scan(
			&totalRecords,
			&like.LikeID,
			&like.UserID,
			&like.EpisodeID,
			&like.CommentText,
			&like.LikeCount,
			&like.CreatedAt,
			&like.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		likes = append(likes, &like)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return likes, metadata, nil
}

func ValidateLike(v *validator.Validator, likeComment *LikeComment) {
	// Check if the title field is empty.
	v.Check(likeComment.CommentText != "", "CommentText", "must be provided")
	// Check if the title field is not more than 100 characters.
	v.Check(len(likeComment.CommentText) <= 100, "CommentText", "must not be more than 100 bytes long")
	// Check if the description field is not more than 1000 characters.
	v.Check(likeComment.LikeCount >= 0, "like_count", "must not be negative")

}
