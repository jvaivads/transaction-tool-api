package summarizer

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
	"transaction-tool-api/src/internal/notifier"

	"github.com/stretchr/testify/assert"
)

func TestServiceResumeTransactions(t *testing.T) {
	var (
		date      = time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC)
		customErr = errors.New("custom error")
		bankTnxs  = transactions{
			items: []transaction{
				{amount: 10, date: date},
			},
			userID: 1,
		}
		msg, _ = Resume{
			Balance:           "10.00",
			CreditAvg:         "10.00",
			DebitAvg:          "0.00",
			MonthTransactions: []MonthTransaction{{time.December.String(), 1}},
		}.ToHTML(resumeHTMLTemplate)
	)

	tests := []struct {
		name         string
		transactions transactions
		mockApplier  func(rm *repositoryMock, nm *notifier.Mock)
		expected     error
	}{
		{
			name:         "no transactions",
			transactions: transactions{},
			expected:     nil,
		},
		{
			name: "init transactional operations fails",
			transactions: transactions{
				items: []transaction{
					{},
				},
			},
			mockApplier: func(rm *repositoryMock, nm *notifier.Mock) {
				rm.On("initTransactionalOperations").Return(nil, customErr).Once()
			},
			expected: fmt.Errorf("error creating repository transaction due to: %w", customErr),
		},
		{
			name:         "error getting user",
			transactions: bankTnxs,
			mockApplier: func(rm *repositoryMock, nm *notifier.Mock) {
				rm.On("initTransactionalOperations").Return(tx{}, nil).Once()
				rm.On("getUserByID", tx{}, int64(1)).Return(nil, customErr).Once()
				rm.On(
					"finishTransactionalOperations",
					tx{}, fmt.Errorf("error getting user due to: %w", customErr)).
					Return(customErr).Once()
			},
			expected: customErr,
		},
		{
			name:         "save bank transactions fails",
			transactions: bankTnxs,
			mockApplier: func(rm *repositoryMock, nm *notifier.Mock) {
				rm.On("initTransactionalOperations").Return(tx{}, nil).Once()
				rm.On("getUserByID", tx{}, int64(1)).Return(User{Email: "email"}, nil).Once()
				rm.On("saveBankTransactions", tx{}, bankTnxs).Return(customErr).Once()
				rm.On(
					"finishTransactionalOperations",
					tx{}, fmt.Errorf("error saving transactions due to: %w", customErr)).
					Return(customErr).Once()
			},
			expected: customErr,
		},
		{
			name:         "notifier user fails",
			transactions: bankTnxs,
			mockApplier: func(rm *repositoryMock, nm *notifier.Mock) {
				rm.On("initTransactionalOperations").Return(tx{}, nil).Once()
				rm.On("getUserByID", tx{}, int64(1)).Return(User{Email: "email"}, nil).Once()
				rm.On("saveBankTransactions", tx{}, bankTnxs).Return(nil).Once()
				nm.On("NotifyToUser", msg, "email").Return(customErr).Once()
				rm.On(
					"finishTransactionalOperations", tx{},
					fmt.Errorf("error notifying transactions to user id %d due to: %w", bankTnxs.userID, customErr)).
					Return(customErr).Once()
			},
			expected: customErr,
		},
		{
			name:         "finish transaction fails",
			transactions: bankTnxs,
			mockApplier: func(rm *repositoryMock, nm *notifier.Mock) {
				rm.On("initTransactionalOperations").Return(tx{}, nil).Once()
				rm.On("saveBankTransactions", tx{}, bankTnxs).Return(nil).Once()
				rm.On("getUserByID", tx{}, int64(1)).Return(User{Email: "email"}, nil).Once()
				nm.On("NotifyToUser", msg, "email").Return(nil).Once()
				rm.On("finishTransactionalOperations", tx{}, error(nil)).Return(customErr).Once()
			},
			expected: customErr,
		},
		{
			name:         "user notified successfully",
			transactions: bankTnxs,
			mockApplier: func(rm *repositoryMock, nm *notifier.Mock) {
				rm.On("initTransactionalOperations").Return(tx{}, nil).Once()
				rm.On("getUserByID", tx{}, int64(1)).Return(User{Email: "email"}, nil).Once()
				rm.On("saveBankTransactions", tx{}, bankTnxs).Return(nil).Once()
				nm.On("NotifyToUser", msg, "email").Return(nil).Once()
				rm.On("finishTransactionalOperations", tx{}, error(nil)).Return(nil).Once()
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			notifierMock := &notifier.Mock{}
			repoMock := &repositoryMock{}
			if test.mockApplier != nil {
				test.mockApplier(repoMock, notifierMock)
				defer repoMock.AssertExpectations(t)
				defer notifierMock.AssertExpectations(t)
			}
			serv := service{
				repository: repoMock,
				notifier:   notifierMock,
			}
			assert.Equal(t, test.expected, serv.notifyResume(context.TODO(), test.transactions))
		})
	}
}
