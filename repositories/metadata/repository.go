package metadata

import (
	"context"

	"github.com/pkg/errors"

	"github.com/teran/linker/models"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	ResolveLink(ctx context.Context, linkID string) (models.Link, error)
}
