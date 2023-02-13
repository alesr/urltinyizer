package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	getShortURLQuery            string = "SELECT short_url FROM urls WHERE long_url = $1"
	getLongURLQuery             string = "SELECT long_url FROM urls WHERE short_url = $1"
	geStatsQuery                string = "SELECT hits FROM urls WHERE short_url = $1"
	updateHitsAndLastHitAtQuery string = "UPDATE urls SET hits = hits + 1, last_hit_at = NOW() WHERE short_url = $1"
	saveShortURLQuery           string = "INSERT INTO urls (short_url, long_url) VALUES ($1, $2)"
)

// DB defines a interface with the methods from sqlx.DB struct.
type db interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// PostgreSQL is a repository that implements the Repository interface.
type PostgreSQL struct {
	logger *zap.Logger
	dbConn db
}

// NewPostgreSQL creates a new PostgreSQL repository.
func NewPostgreSQL(logger *zap.Logger, dbConn db) *PostgreSQL {
	return &PostgreSQL{
		logger: logger,
		dbConn: dbConn,
	}
}

// GetShortURL returns the short URL for a given long URL.
func (p *PostgreSQL) GetShortURL(ctx context.Context, longURL string) (string, error) {
	var shortURL string
	if err := p.dbConn.GetContext(ctx, &shortURL, getShortURLQuery, longURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("could not get short URL from database: %w", err)
	}
	return shortURL, nil
}

// GetLongURL returns the long URL for a given short URL.
func (p *PostgreSQL) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	tx, err := p.dbConn.BeginTxx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	var longURL string
	if err := tx.GetContext(ctx, &longURL, getLongURLQuery, shortURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("could not get long URL from database: %w", err)
	}

	if _, err := tx.ExecContext(ctx, updateHitsAndLastHitAtQuery, shortURL); err != nil {
		return "", fmt.Errorf("could not update hits and last_hit_at: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("could not commit transaction: %w", err)
	}
	return longURL, nil
}

// SaveShortURL saves a short URL to the database.
func (p *PostgreSQL) SaveShortURL(ctx context.Context, shortURL, longURL string) error {
	if _, err := p.dbConn.ExecContext(ctx, saveShortURLQuery, shortURL, longURL); err != nil {
		return fmt.Errorf("could not save short URL to database: %w", err)
	}
	return nil
}

// Get sats for a given short URL.
func (p *PostgreSQL) GetStats(ctx context.Context, shortURL string) (int, error) {
	var hits int
	if err := p.dbConn.GetContext(ctx, &hits, geStatsQuery, shortURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("could not get hits from database: %w", err)
	}
	return hits, nil
}
