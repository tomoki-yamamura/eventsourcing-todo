package repository

import "context"

type Transaction interface {
	RWTx(ctx context.Context, fn func(ctx context.Context) error) error
	AfterCommit(fn func() error)
}
