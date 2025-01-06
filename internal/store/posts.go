package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

// Post model
type Post struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	UserID    int64    `json:"user_id"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetPostById(ctx context.Context, postID int64) (*Post, error) {
	query := `SELECT * FROM posts WHERE id=$1`
	row := s.db.QueryRow(query, postID)

	var post Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags), // notice: need to convert back to pq.Array
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &post, nil
}
