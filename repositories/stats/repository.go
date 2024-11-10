package stats

import (
	"context"

	"github.com/teran/linker/models"
)

type Repository interface {
	LogRequest(ctx context.Context, req models.Request) error
}
