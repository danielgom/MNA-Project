package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dBContextTimeout = 10

// DBConn interface is used to call the DB.
type DBConn interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// FakeScanner Probably won't use but should be used for mocking a scanner.
type FakeScanner interface {
	Scan(dest ...interface{}) error
}

// InitDatabase returns a pool from DB configuration.
func InitDatabase(c *Config) *pgxpool.Pool {
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s%s/%s", c.DB.User,
		c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name)

	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Panicln("could not connect to database,", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*dBContextTimeout)
	defer cancel()

	err = conn.Ping(ctx)
	if err != nil {
		log.Panicln("could not connect to database,", err.Error())
	}

	return conn
}
