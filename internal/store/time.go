package store

import (
	"database/sql"
	"time"
)

type DBTimeStore struct {
	db *sql.DB
}

func parsePostgresTime(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}

func toPostgresTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func (s *DBTimeStore) GetDBTimeZone() (string, error) {
	query := `
		SHOW TIME ZONE
	`

	var dbTimeZone string

	err := s.db.QueryRow(query).Scan(&dbTimeZone)
	if err != nil {
		return "", err
	}
	return dbTimeZone, nil
}
