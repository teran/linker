package service

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/teran/linker/models"
)

var _ Service = (*Mock)(nil)

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Redirect(ctx context.Context, req models.Request) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}
