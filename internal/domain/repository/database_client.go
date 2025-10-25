package repository

import "github.com/jmoiron/sqlx"

type DatabaseClient interface {
	GetDB() *sqlx.DB
	Close() error
}
