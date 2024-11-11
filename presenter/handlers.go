package presenter

import (
	"net/http"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"

	"github.com/teran/linker/models"
	"github.com/teran/linker/service"
)

type Config struct {
	Domain             string
	CookieName         string
	URLService         service.Service
	CookieIDGenerator  func() string
	TimeNowProvider    func() time.Time
	ClientIPHeaderName string
}

type Handlers interface {
	processRedirect(echo.Context) error

	Register(e *echo.Echo)
}

type handlers struct {
	cfg *Config
}

func New(cfg *Config) Handlers {
	return &handlers{
		cfg: cfg,
	}
}

func (h *handlers) processRedirect(c echo.Context) error {
	cookie, err := c.Cookie(h.cfg.CookieName)
	if err != nil {
		cookie = &http.Cookie{
			Name:     h.cfg.CookieName,
			Path:     "/",
			Domain:   h.cfg.Domain,
			Expires:  h.cfg.TimeNowProvider().Add(24 * 730 * time.Hour),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Value:    h.cfg.CookieIDGenerator(),
		}

		c.SetCookie(cookie)
	}

	clientIP := strings.SplitN(c.Request().RemoteAddr, ":", 2)[0]
	if h.cfg.ClientIPHeaderName != "" {
		clientIP = c.Request().Header.Get(h.cfg.ClientIPHeaderName)
	}

	url, err := h.cfg.URLService.Redirect(c.Request().Context(), models.Request{
		Timestamp:  h.cfg.TimeNowProvider().UTC(),
		LinkID:     c.Param("linkID"),
		ClientIP:   clientIP,
		CookieID:   cookie.Value,
		UserAgent:  c.Request().UserAgent(),
		Parameters: models.Parameters(c.Request().URL.Query()),
	})
	if err != nil {
		return mapServiceErrors(c, err)
	}

	c.Response().Header().Set("Server", "linker/1.0")
	c.Response().Header().Set("Referrer-Policy", "unsafe-url")
	c.Response().Header().Set("Content-Security-Policy", "referrer always")

	return c.Redirect(http.StatusFound, url)
}

func (h *handlers) robotsTxt(c echo.Context) error {
	return c.Blob(http.StatusOK, "text/plain", []byte("User-agent: *\nDisallow: /\n"))
}

func (h *handlers) Register(e *echo.Echo) {
	e.GET("/robots.txt", h.robotsTxt)
	e.GET("/:linkID", h.processRedirect)
}
