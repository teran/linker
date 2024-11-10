package presenter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/teran/linker/models"
	"github.com/teran/linker/service"
)

const (
	testDomain     = "test-domain.example.com"
	testCookieName = "ccID"
)

func (s *handlersTestSuite) TestFreshClientRedirect() {
	s.svcM.On("Redirect", models.Request{
		LinkID:     "linkIDStub",
		ClientIP:   "127.0.0.1",
		CookieID:   "test-cookie-id",
		UserAgent:  "test-client/1.0",
		Parameters: map[string][]string{},
	}).Return("https://test-redirect-url.example.com/", nil).Once()

	req, err := http.NewRequest(http.MethodGet, s.srv.URL+"/linkIDStub", nil)
	s.Require().NoError(err)

	req.Header.Set("User-Agent", "test-client/1.0")

	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := c.Do(req)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusFound, resp.StatusCode)
	s.Require().Equal("https://test-redirect-url.example.com/", resp.Header.Get("Location"))
	s.Require().Equal(
		"ccID=test-cookie-id; Path=/; Domain=test-domain.example.com; Expires=Mon, 09 Nov 2026 01:02:03 GMT; HttpOnly; Secure; SameSite=Strict",
		resp.Header.Get("Set-Cookie"),
	)
}

func (s *handlersTestSuite) TestFreshClientRedirectNotFound() {
	s.svcM.On("Redirect", models.Request{
		LinkID:     "linkIDStub",
		ClientIP:   "127.0.0.1",
		CookieID:   "test-cookie-id",
		UserAgent:  "test-client/1.0",
		Parameters: map[string][]string{},
	}).Return("", service.ErrNotFound).Once()

	req, err := http.NewRequest(http.MethodGet, s.srv.URL+"/linkIDStub", nil)
	s.Require().NoError(err)

	req.Header.Set("User-Agent", "test-client/1.0")

	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := c.Do(req)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

// Definitions ...

type handlersTestSuite struct {
	suite.Suite

	srv  *httptest.Server
	svcM *service.Mock
}

func (s *handlersTestSuite) SetupTest() {
	s.svcM = service.NewMock()

	hndlr := New(&Config{
		Domain:            testDomain,
		CookieName:        testCookieName,
		URLService:        s.svcM,
		CookieIDGenerator: func() string { return "test-cookie-id" },
		TimeNowProvider:   func() time.Time { return time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC).UTC() },
	})

	e := echo.New()
	hndlr.Register(e)

	s.srv = httptest.NewServer(e)
}

func (s *handlersTestSuite) TearDownTest() {
	s.srv.Close()

	s.svcM.AssertExpectations(s.T())
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, &handlersTestSuite{})
}
