package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates a new SQLite storage instance, connecting to the database at the
// given storagePath. If the database does not exist, it will be created. If the
// database exists, but the table does not, it will be created. If the database
// exists and the table exists, no changes will be made to the database.
//
// The returned Storage instance is ready to use.
func New(storagePath string) (*Storage, error) {
	const operation = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", operation, err)
	}

	return &Storage{db: db}, nil
}

// SaveURL saves the given URL with the given alias in the database.
// If the alias is empty, it generates a random one.
// If the URL already exists in the database, it returns ErrURLExists.
// Otherwise, it returns the ID of the newly created URL.
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const operation = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", operation, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last inserted ID: %w", operation, err)
	}

	return id, nil
}

// GetURL retrieves the URL associated with the given alias from the database.
// Returns the URL as a string if found, otherwise returns an error if the alias
// does not exist or if there is an issue executing the query.
func (s *Storage) GetURL(alias string) (string, error) {
	const operation = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", operation, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", operation, err)
	}

	return resURL, nil
}

// DeleteURL deletes a URL from the database
func (s *Storage) DeleteURL(alias string) error {
	const operation = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Exec("DELETE FROM url WHERE alias = ?", alias)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	rowsAffected, err := stmt.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("url not found")
	}

	return nil
}
