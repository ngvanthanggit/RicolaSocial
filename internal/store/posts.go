package store

import (
	"context"
	"database/sql"
	"log"

	"github.com/lib/pq"
)

// Post model
type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	UserID    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentsCount int64 `json:"comments_count"`
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, postID)

	var post Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags), // notice: need to convert back to pq.Array
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `
		DELETE FROM posts WHERE posts.id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return ErrNotFound
	}
	log.Printf("%d row(s) affected", row)
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, updated_at = NOW(), version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, feedQuery PaginatedFeedQuery) ([]PostWithMetaData, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at, p.version, p.tags, COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.user_id = p.user_id
		LEFT JOIN followers f ON f.user_id = p.user_id OR p.user_id = $1
		JOIN users u ON p.user_id = u.id
		WHERE (f.follower_id = $1 OR p.user_id = $1)
		AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
		AND (p.tags @> $5 OR $5 = '{}')
		AND (p.created_at >= $6 AND p.created_at <= $7)
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + feedQuery.Sort + `
		LIMIT $2 OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userID,
		feedQuery.Limit,
		feedQuery.Offset,
		feedQuery.Search,
		pq.Array(feedQuery.Tags),
		feedQuery.Since,
		feedQuery.Until,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostWithMetaData
	for rows.Next() {
		var post PostWithMetaData
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.User.Username,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		post.User.ID = post.UserID
		feed = append(feed, post)
	}

	return feed, nil
}
