//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/alesr/urltinyizer/helper"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetShortURL(t *testing.T) {
	db := helper.SetupDB(t, "../../migrations")
	defer helper.TeardownDB(t, db)

	repo := NewPostgreSQL(zap.NewNop(), db)

	// Insert a new entry into the database.
	_, err := db.Exec("INSERT INTO urls (short_url, long_url) VALUES ('abc', 'https://www.foo.bar')")
	require.NoError(t, err)

	// Get the short URL for the entry we just inserted.
	shortURL, err := repo.GetShortURL(context.TODO(), "https://www.foo.bar")
	require.NoError(t, err)

	require.Equal(t, "abc", shortURL)
}

func TestGetLongURL(t *testing.T) {
	db := helper.SetupDB(t, "../../migrations")
	defer helper.TeardownDB(t, db)

	repo := NewPostgreSQL(zap.NewNop(), db)

	// Insert a new entry into the database.
	_, err := db.Exec("INSERT INTO urls (short_url, long_url) VALUES ('abc', 'https://www.foo.bar')")
	require.NoError(t, err)

	// Get the long URL for the entry we just inserted.
	longURL, err := repo.GetLongURL(context.TODO(), "abc")
	require.NoError(t, err)

	require.Equal(t, "https://www.foo.bar", longURL)
}

func TestSaveShortURL(t *testing.T) {
	db := helper.SetupDB(t, "../../migrations")
	defer helper.TeardownDB(t, db)

	repo := NewPostgreSQL(zap.NewNop(), db)

	// Save a new short URL.
	err := repo.SaveShortURL(context.TODO(), "abc", "https://www.foo.bar")
	require.NoError(t, err)

	// Get the long URL for the entry we just inserted.
	longURL, err := repo.GetLongURL(context.TODO(), "abc")
	require.NoError(t, err)

	require.Equal(t, "https://www.foo.bar", longURL)
}

func TestGetStats(t *testing.T) {
	db := helper.SetupDB(t, "../../migrations")
	defer helper.TeardownDB(t, db)

	repo := NewPostgreSQL(zap.NewNop(), db)

	// Insert a new entry into the database.
	_, err := db.Exec("INSERT INTO urls (short_url, long_url) VALUES ('abc', 'https://www.foo.bar')")
	require.NoError(t, err)

	// Get the long URL for the entry we just inserted so the hit count is 1.
	_, err = repo.GetLongURL(context.TODO(), "abc")
	require.NoError(t, err)

	// Get the stats for the entry we just inserted.
	stats, err := repo.GetStats(context.TODO(), "abc")
	require.NoError(t, err)

	require.Equal(t, 1, stats)
}
