package summarizer

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) initTransactionalOperations(_ context.Context) (tx, error) {
	var (
		txn  tx
		args = m.Called()
	)

	if value, ok := args.Get(0).(tx); ok {
		txn = value
	}
	return txn, args.Error(1)
}

func (m *repositoryMock) finishTransactionalOperations(_ context.Context, txn tx, error error) error {
	args := m.Called(txn, error)
	return args.Error(0)
}

func (m *repositoryMock) saveBankTransactions(_ context.Context, txn tx, bankTxns transactions) error {
	args := m.Called(txn, bankTxns)
	return args.Error(0)
}

func (m *repositoryMock) getUserByID(_ context.Context, txn tx, userID int64) (User, error) {
	var (
		user User
		args = m.Called(txn, userID)
	)

	if value, ok := args.Get(0).(User); ok {
		user = value
	}

	return user, args.Error(1)
}
