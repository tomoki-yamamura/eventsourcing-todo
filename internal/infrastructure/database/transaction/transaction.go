package transaction

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
)

type txKeyType string

const TxKey txKeyType = "tx"

type transaction struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) repository.Transaction {
	return &transaction{db: db}
}

func (t *transaction) RWTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.runTx(ctx, sql.LevelRepeatableRead, fn)
}

func (t *transaction) runTx(ctx context.Context, level sql.IsolationLevel, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{Isolation: level})
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	var committed bool
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := fn(ctxWithTx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	committed = true
	return nil
}

func GetTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(TxKey).(*sqlx.Tx)
	if !ok || tx == nil {
		return nil, errors.New("transaction not found in context")
	}
	return tx, nil
}

func WithTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
