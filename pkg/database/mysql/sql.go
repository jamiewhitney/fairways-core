package database_sql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jamiewhitney/fairways-core/pkg/database"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
)

type DB struct {
	Pool *sql.DB
}

func NewFromConfig(cfg *database.Config) (*DB, error) {
	var dsn string

	dsn = fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true", cfg.UserPass, cfg.Host, cfg.Port, cfg.Name)

	client, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(); err != nil {
		return nil, err
	}

	return &DB{Pool: client}, err

}

func (db *DB) Ping() error {
	return db.Pool.Ping()
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) InTx(ctx context.Context, f func(tx *sql.Tx) error) error {
	tx, err := db.Pool.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: false,
	})
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if err := f(tx); err != nil {
		logging.FromContext(ctx).Error(err)
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("rolling back transaction: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
