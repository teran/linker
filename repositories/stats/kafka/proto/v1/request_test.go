package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/teran/linker/models"
)

func TestRequestMapping(t *testing.T) {
	r := require.New(t)

	req := models.Request{
		Timestamp: time.Date(2024, 11, 3, 1, 2, 3, 0, time.UTC),
		LinkID:    "test-link-id",
		ClientIP:  "127.0.0.1",
		CookieID:  "test-cookie-id",
		UserAgent: "test-user-agent/1.0",
		Parameters: models.Parameters{
			"key1": []string{"value1"},
			"key2": []string{"value2"},
		},
		Referrer: "test-referer",
	}

	pr, err := NewRequest(req)
	r.NoError(err)

	rReq, err := pr.ToSvc()
	r.NoError(err)
	r.Equal(req, rReq)
}
