package mock

import (
	"context"

	"github.com/teran/linker/models"
	"github.com/teran/linker/repositories/metadata/mock"
)

type Mock struct {
	mock.Mock
}

func New() *Mock {
	return &Mock{}
}

func (m *Mock) LogRequest(_ context.Context, req models.Request) error {
	args := m.Called(req)
	return args.Error(0)
}
