package store

import (
	"database/sql"
	"time"
)

var (
	DefaultStartTime = time.Date(2003, time.January, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
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

func (s *DBTimeStore) DBTimeSetup() error {
	dbTimeZone, err := s.GetDBTimeZone()
	if err != nil {
		return err
	}
	location, err := time.LoadLocation(dbTimeZone)
	if err != nil {
		return err
	}
	time.Local = location

	return nil
}

func TimeNowString() string {
	return time.Now().Format(time.RFC3339)
}
