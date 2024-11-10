package presenter

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/teran/linker/service"
)

func mapServiceErrors(c echo.Context, err error) error {
	if errors.Is(err, service.ErrNotFound) {
		return c.Blob(http.StatusNotFound, "text/plain", []byte("Not Found"))
	}

	return err
}
