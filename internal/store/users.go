package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("the email already exists")
	ErrDuplicateUsername = errors.New("the username already exists")
	ErrTokenExpired      = errors.New("the token has already expired")
	ErrUserActivated     = errors.New("the user has already been activated")
)

// User model
type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.hash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) GetUserById(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, userID)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// this function will execute the queries and return error if there's is
		// something wrong

		// create
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// we will change the create() so that we can pass the transaction
		// the problem now is that where to create the transaction and manage all the
		// executions inside a transaction
		// where?
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		// invite

		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	// as this includes multiple steps of updating db -> transaction
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// - find the user with this token
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		// - update the user (active status)
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}
		// - remove the token from db
		if err := s.deleteInvitation(ctx, tx, token); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int64) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT i.user_id, u.username, u.email, u.created_at, u.is_active, i.expiry
		FROM user_invitations i
		JOIN users u ON i.user_id = u.id
		WHERE token = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	row := tx.QueryRowContext(ctx, query, hashToken)
	var user User
	var queryExp time.Time
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
		&queryExp,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	// user is activated
	if user.IsActive {
		return nil, ErrUserActivated
	}
	// the token is expired
	if time.Now().After(queryExp) {
		return nil, ErrTokenExpired
	}

	return &user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users
		SET username=$2, email=$3, created_at=$4, is_active=$5
		WHERE id=$1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.CreatedAt, user.IsActive)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM users WHERE id=$1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UserStore) deleteInvitation(ctx context.Context, tx *sql.Tx, token string) error {
	query := `
		DELETE FROM user_invitations WHERE token=$1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	res, err := tx.ExecContext(ctx, query, hashToken)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM user_invitations WHERE user_id=$1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
