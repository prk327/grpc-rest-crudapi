package database

import (
	"context"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	_ "github.com/vertica/vertica-sql-go"

	"github.com/prk327/grpc-rest-crudapi/server/config"
)

func DatabaseConfig() config.DatabaseConfig {
	return config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5433",
		User:     "dbadmin",
		Password: "admin1234",
		Name:     "vdb",
		Schema:   "omniq",
	}
}

func TestCRUDOperations(t *testing.T) {
	db, err := New(DatabaseConfig())
	require.NoError(t, err)
	defer db.Close()

	// Run migrations
	err = db.Migrate(context.Background())
	require.NoError(t, err)

	// Run tests...

	// Optional: Rollback changes
	err = goose.DownTo(db.Conn, "migrations", 0)
	require.NoError(t, err)
}
