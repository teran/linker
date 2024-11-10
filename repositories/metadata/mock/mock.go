package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/teran/linker/models"
	"github.com/teran/linker/repositories/metadata"
)

var _ metadata.Repository = (*Mock)(nil)

type Mock struct {
	mock.Mock
}

func New() *Mock {
	return &Mock{}
}

func (m *Mock) ResolveLink(ctx context.Context, linkID string) (models.Link, error) {
	args := m.Called(linkID)
	return args.Get(0).(models.Link), args.Error(1)
}
