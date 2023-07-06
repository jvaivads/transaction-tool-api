package summarizer

import (
	"context"
	"database/sql"
)

type tx struct {
	client *sql.Tx
}

type repository interface {
	initTransactionalOperations(context.Context) (tx, error)
	finishTransactionalOperations(context.Context, tx, error) error
	saveBankTransactions(context.Context, tx, transactions) error
}
