package summarizer

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (m *serviceMock) notifyResume(_ context.Context, txns transactions) (err error) {
	args := m.Called(txns)
	return args.Error(0)
}
