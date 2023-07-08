package notifier

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gopkg.in/mail.v2"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) NotifyToUser(_ context.Context, msg string, userID int64) error {
	args := m.Called(msg, userID)
	return args.Error(0)
}

type dialerMock struct {
	mock.Mock
}

func (m *dialerMock) DialAndSend(messages ...*mail.Message) error {
	args := m.Called(messages)
	return args.Error(0)
}
