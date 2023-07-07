package summarizer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type repository interface {
	initTransactionalOperations(context.Context) (tx, error)
	finishTransactionalOperations(context.Context, tx, error) error
	saveBankTransactions(context.Context, tx, transactions) error
}

type tx struct {
	client *sql.Tx
}

func (t tx) Exec(ctx context.Context, query string, params []any) (sql.Result, error) {
	if t.client == nil {
		return nil, errors.New("transaction client is nil")
	}
	return t.client.ExecContext(ctx, query, params...)
}

type sqlRepository struct {
	client *sql.DB
}

func (r sqlRepository) initTransactionalOperations(ctx context.Context) (tx, error) {
	tnx, err := r.client.BeginTx(ctx, nil)
	if err != nil {
		return tx{}, err
	}
	return tx{client: tnx}, nil
}

func (r sqlRepository) finishTransactionalOperations(_ context.Context, tnx tx, err error) error {
	if err != nil {
		_ = tnx.client.Rollback()
		return err
	}
	return tnx.client.Commit()
}

func (r sqlRepository) saveBankTransactions(ctx context.Context, tnx tx, bankTxns transactions) error {
	if bankTxns.items == nil {
		return nil
	}
	var (
		query              = `INSERT INTO transaction (user_id, amount, date_created) VALUES %s`
		transactionFormat  = `(?,?,?)`
		params             = make([]any, 0, 3*len(bankTxns.items))
		transactionFormats = make([]string, 0, len(bankTxns.items))
	)

	for _, bankTxn := range bankTxns.items {
		transactionFormats = append(transactionFormats, transactionFormat)
		params = append(params, bankTxns.userID, bankTxn.amount, bankTxn.date)
	}

	query = fmt.Sprintf(query, strings.Join(transactionFormats, ","))

	result, err := tnx.Exec(ctx, query, params)
	if err != nil {
		return fmt.Errorf("error inserting transactions due to: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting total inserted transactions due to: %w", err)
	}

	if rowsAffected != int64(len(bankTxns.items)) {
		return fmt.Errorf(
			"total affected rows (%d) mismatch with total transactions (%d)", rowsAffected, len(bankTxns.items))
	}

	return nil
}
