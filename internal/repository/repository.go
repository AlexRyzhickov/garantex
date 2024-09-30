package repository

import (
	"context"
	"garantex/internal/models"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func New(db *pgxpool.Pool) *DB {
	return &DB{db}
}

func BuildUpsertQuery(table string, price models.Price) string {
	buf := make([]byte, 0)
	buf = append(buf, "INSERT INTO "...)
	buf = append(buf, table...)
	buf = append(buf, " (ts, ask_price, bid_price) VALUES "...)
	buf = append(buf, "("...)
	buf = append(buf, strconv.Itoa(int(price.Timestamp))...)
	buf = append(buf, ',')
	buf = strconv.AppendFloat(buf, price.AskPrice, 'f', -1, 32)
	buf = append(buf, ',')
	buf = strconv.AppendFloat(buf, price.BidPrice, 'f', -1, 32)
	buf = append(buf, ")"...)
	buf = append(buf, " ON CONFLICT (ts) DO UPDATE SET ts = EXCLUDED.ts, ask_price = EXCLUDED.ask_price, bid_price = EXCLUDED.bid_price;"...)
	return string(buf)
}

func (db *DB) Upsert(price models.Price) error {
	tx, err := db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := BuildUpsertQuery("prices", price)

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
