package helper

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

const (
	postgresDriverName string = "postgres"
	dbHost             string = "localhost"
	dbPort             string = "5432"
	dbUser             string = "user"
	dbPass             string = "password"
	dbName             string = "urltinyizer"
)

func SetupDB(t *testing.T, migrationsDir string) *sqlx.DB {
	db, err := sqlx.Open(postgresDriverName, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName),
	)
	require.NoError(t, err)

	require.NoError(t, goose.Up(db.DB, migrationsDir))
	return db
}

func TeardownDB(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec("TRUNCATE TABLE urls;")
	require.NoError(t, err)

	require.NoError(t, db.Close())
}
