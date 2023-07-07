package summarizer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDBMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return db, mock
}

func TestSQLRepositorySaveBankTransactions(t *testing.T) {
	var (
		customErr = errors.New("custom error")
		date      = time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)
		bankTxns  = transactions{
			items: []transaction{
				{
					amount: 10,
					date:   date,
				},
				{
					amount: 10,
					date:   date.Add(time.Hour),
				},
			},
			userID: 5,
		}
		query  = regexp.QuoteMeta(`INSERT INTO transaction (user_id, amount, date_created) VALUES (?,?,?),(?,?,?)`)
		params = []driver.Value{
			bankTxns.userID, bankTxns.items[0].amount, bankTxns.items[0].date,
			bankTxns.userID, bankTxns.items[1].amount, bankTxns.items[1].date,
		}
	)
	tests := []struct {
		name        string
		bankTxns    transactions
		mockApplier func(sqlmock.Sqlmock)
		expected    error
	}{
		{
			name:        "no transactions",
			bankTxns:    transactions{},
			mockApplier: nil,
			expected:    nil,
		},
		{
			name:     "error executing query",
			bankTxns: bankTxns,
			mockApplier: func(m sqlmock.Sqlmock) {
				m.ExpectExec(query).WithArgs(params...).WillReturnError(customErr)
			},
			expected: fmt.Errorf("error inserting transactions due to: %w", customErr),
		},
		{
			name:     "error getting rows affected",
			bankTxns: bankTxns,
			mockApplier: func(m sqlmock.Sqlmock) {
				m.ExpectExec(query).WithArgs(params...).WillReturnResult(sqlmock.NewErrorResult(customErr))
			},
			expected: fmt.Errorf("error getting total inserted transactions due to: %w", customErr),
		},
		{
			name:     "error total affected rows mismatch",
			bankTxns: bankTxns,
			mockApplier: func(m sqlmock.Sqlmock) {
				m.ExpectExec(query).WithArgs(params...).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expected: fmt.Errorf(
				"total affected rows (%d) mismatch with total transactions (%d)", 1, 2),
		},
		{
			name:     "transactions inserted successfully",
			bankTxns: bankTxns,
			mockApplier: func(m sqlmock.Sqlmock) {
				m.ExpectExec(query).WithArgs(params...).WillReturnResult(sqlmock.NewResult(0, 2))
			},
			expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := newDBMock()
			mock.ExpectBegin()
			if test.mockApplier != nil {
				test.mockApplier(mock)
			}
			tnx, err := db.Begin()
			if err != nil {
				require.Nil(t, err)
			}
			repo := sqlRepository{client: db}
			assert.Equal(t, test.expected, repo.saveBankTransactions(context.TODO(), tx{tnx}, test.bankTxns))

		})
	}
}
