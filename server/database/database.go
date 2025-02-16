package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/vertica/vertica-sql-go"

	"github.com/prk327/grpc-rest-crudapi/server/config"
)

type Database struct {
	Conn *sql.DB
}

func New(cfg config.DatabaseConfig) (*Database, error) {
	connStr := fmt.Sprintf(
		"vertica://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	println(connStr)

	db, err := sql.Open("vertica", connStr)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &Database{Conn: db}, nil
}

func (d *Database) Close() error {
	return d.Conn.Close()
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Conn.PingContext(ctx)
}

func (d *Database) ValidateSchema(ctx context.Context, schema string) error {
	var exists bool
	err := d.Conn.QueryRowContext(ctx,
		`SELECT EXISTS(
            SELECT schema_name 
            FROM v_catalog.schemata 
            WHERE schema_name = $1
        )`, schema).Scan(&exists)

	if err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	if !exists {
		return fmt.Errorf("schema '%s' does not exist", schema)
	}

	return nil
}
