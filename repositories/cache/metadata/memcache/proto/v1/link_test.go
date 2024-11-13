package v1

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teran/linker/models"
)

func TestLinkMappin(t *testing.T) {
	r := require.New(t)

	l := models.Link{
		DestinationURL: "https://test-url.example.com",
		Parameters: models.Parameters{
			"key1": []string{"value1", "value2"},
			"key2": []string{"value3"},
		},
	}

	pl, err := NewLink(l)
	r.NoError(err)

	rl, err := pl.ToSvc()
	r.NoError(err)

	r.Equal(l, rl)
}
