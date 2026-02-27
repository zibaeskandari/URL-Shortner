package db

import (
	"URLShortner/internal/infrastructure/persistence/ent"
	"URLShortner/pkg"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDb struct {
	config pkg.DatabaseConfig
	db     *sql.DB
	driver dialect.Driver
	pool   *pgxpool.Pool
	client *ent.Client
}

func NewPostgresDb(config pkg.DatabaseConfig) (*PostgresDb, error) {
	// TODO: Add validation here
	return &PostgresDb{
		config: config,
		db:     nil,
		driver: nil,
		client: nil,
	}, nil
}

func (pgDb *PostgresDb) Open() (*ent.Client, error) {
	ctx := context.Background()
	source := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", pgDb.config.Username, pgDb.config.Password, pgDb.config.Host, pgDb.config.Port, pgDb.config.Database, pgDb.config.SSLMode)
	poolConf, err := pgxpool.ParseConfig(source)
	if err != nil {
		log.Fatalf("failed pars db config: %v", err)
	}

	poolConf.MaxConns = pgDb.config.MaxConns
	poolConf.MinConns = pgDb.config.MinConns
	poolConf.MaxConnLifetime = time.Duration(pgDb.config.MaxConnLifetime) * time.Minute
	poolConf.MaxConnIdleTime = time.Duration(pgDb.config.MaxConnIdleTime) * time.Minute

	pgDb.pool, err = pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		log.Fatalf("failed to create connection pool: %v", err)
	}

	pgDb.db = stdlib.OpenDBFromPool(pgDb.pool)

	pgDb.driver = entsql.OpenDB(dialect.Postgres, pgDb.db)
	pgDb.client = ent.NewClient(ent.Driver(pgDb.driver))
	return pgDb.client, nil
}

func (pgDb *PostgresDb) Close() error {
	if pgDb.client == nil {
		return nil
	}
	return pgDb.client.Close()
}
