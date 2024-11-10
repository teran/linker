package service

import (
	"github.com/pkg/errors"

	"github.com/teran/linker/repositories/metadata"
)

var ErrNotFound = errors.New("not found")

func mapRepositoryErrors(err error) error {
	if errors.Is(err, metadata.ErrNotFound) {
		return ErrNotFound
	}

	return err
}
