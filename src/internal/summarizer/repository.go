package summarizer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func NewRepository(client *sql.DB) Repository {
	return repository{client: client}
}

type Repository interface {
	initTransactionalOperations(context.Context) (tx, error)
	finishTransactionalOperations(context.Context, tx, error) error
	saveBankTransactions(context.Context, tx, transactions) error
	getUserByID(context.Context, tx, int64) (User, error)
}

type tx struct {
	client *sql.Tx
}

func (t tx) Exec(ctx context.Context, query string, params ...any) (sql.Result, error) {
	if t.client == nil {
		return nil, errors.New("transaction client is nil")
	}
	return t.client.ExecContext(ctx, query, params...)
}

func (t tx) QueryRow(ctx context.Context, query string, params ...any) (*sql.Row, error) {
	if t.client == nil {
		return nil, errors.New("transaction client is nil")
	}
	return t.client.QueryRowContext(ctx, query, params...), nil
}

type repository struct {
	client *sql.DB
}

func (r repository) initTransactionalOperations(ctx context.Context) (tx, error) {
	tnx, err := r.client.BeginTx(ctx, nil)
	if err != nil {
		return tx{}, err
	}
	return tx{client: tnx}, nil
}

func (r repository) finishTransactionalOperations(_ context.Context, tnx tx, err error) error {
	if err != nil {
		_ = tnx.client.Rollback()
		return err
	}
	return tnx.client.Commit()
}

func (r repository) saveBankTransactions(ctx context.Context, tnx tx, bankTxns transactions) error {
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

	result, err := tnx.Exec(ctx, query, params...)
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

type User struct {
	userID int64
	name   string
	email  string
}

func (r repository) getUserByID(ctx context.Context, tnx tx, userID int64) (User, error) {
	var (
		user  User
		query = `SELECT id, name, email FROM user WHERE id = ?`
	)

	row, err := tnx.QueryRow(ctx, query, userID)
	if err != nil {
		return User{}, err
	}

	if err = row.Scan(&user.userID, &user.name, &user.email); err != nil {
		return User{}, fmt.Errorf("error scanning user by id %d due to: %w", userID, err)
	}

	return user, nil
}
