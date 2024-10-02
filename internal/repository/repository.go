package repository

import (
	"context"
	"garantex/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func New(db *pgxpool.Pool) *DB {
	return &DB{db}
}

func (db *DB) Upsert(price models.Price) error {
	tx, err := db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := upsertPrice(price)

	_, err = tx.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
