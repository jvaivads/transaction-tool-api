package notifier

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) NotifyToUser(_ context.Context, msg string, userID int64) error {
	args := m.Called(msg, userID)
	return args.Error(0)
}
