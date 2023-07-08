package notifier

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientNotifyToUser(t *testing.T) {
	customErr := errors.New("custom error")

	tests := []struct {
		name        string
		mockApplier func(m *dialerMock)
		expected    error
	}{
		{
			name: "return error",
			mockApplier: func(m *dialerMock) {
				m.On("DialAndSend", mock.Anything).Return(customErr).Once()
			},
			expected: fmt.Errorf("unexpected error sending mail to usier id %d due to: %w", 1, customErr),
		},
		{
			name: "no error",
			mockApplier: func(m *dialerMock) {
				m.On("DialAndSend", mock.Anything).Return(nil).Once()
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dMock := &dialerMock{}
			test.mockApplier(dMock)
			defer dMock.AssertExpectations(t)
			c := client{
				sender: "sender",
				dialer: dMock,
			}
			assert.Equal(t, test.expected, c.NotifyToUser(context.TODO(), "message", 1))
		})
	}
}
