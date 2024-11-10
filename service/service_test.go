package service

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	"github.com/teran/linker/models"
	mdRepoMock "github.com/teran/linker/repositories/metadata/mock"
	statsRepoMock "github.com/teran/linker/repositories/stats/mock"
)

func (s *serviceTestSuite) TestRedirectNoParams() {
	s.mdRepoM.On("ResolveLink", "linkID").Return(models.Link{
		DestinationURL: "https://test-url.example.com",
	}, nil).Once()

	s.statsRepoM.On("LogRequest", models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
	}).Return(nil).Once()

	link, err := s.svc.Redirect(s.ctx, models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
	})
	s.Require().NoError(err)
	s.Require().Equal("https://test-url.example.com", link)
}

func (s *serviceTestSuite) TestRedirectLinkDefaultParams() {
	s.mdRepoM.On("ResolveLink", "linkID").Return(models.Link{
		DestinationURL: "https://test-url.example.com",
		Parameters: models.Parameters{
			"param1": []string{"value1"},
			"param2": []string{"value2", "value3"},
		},
	}, nil).Once()

	s.statsRepoM.On("LogRequest", models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
		Parameters: models.Parameters{
			"param1": []string{"value1"},
			"param2": []string{"value2", "value3"},
		},
	}).Return(nil).Once()

	link, err := s.svc.Redirect(s.ctx, models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
	})
	s.Require().NoError(err)
	s.Require().Equal(
		"https://test-url.example.com?param1=value1&param2=value2&param2=value3",
		link,
	)
}

func (s *serviceTestSuite) TestRedirectLinkDefaultParamsAndRequestParamsNoOverride() {
	s.mdRepoM.On("ResolveLink", "linkID").Return(models.Link{
		DestinationURL: "https://test-url.example.com",
		Parameters: models.Parameters{
			"param1": []string{"value1"},
			"param2": []string{"value2", "value3"},
		},
	}, nil).Once()

	s.statsRepoM.On("LogRequest", models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
		Parameters: models.Parameters{
			"reqParam1": []string{"reqValue1"},
			"reqParam2": []string{"reqValue2", "reqValue3"},
		},
	}).Return(nil).Once()

	link, err := s.svc.Redirect(s.ctx, models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
		Parameters: models.Parameters{
			"reqParam1": []string{"reqValue1"},
			"reqParam2": []string{"reqValue2", "reqValue3"},
		},
	})
	s.Require().NoError(err)
	s.Require().Equal(
		"https://test-url.example.com?param1=value1&param2=value2&param2=value3",
		link,
	)
}

func (s *serviceTestSuite) TestRedirectLinkDefaultParamsAndRequestParamsWithOverride() {
	s.mdRepoM.On("ResolveLink", "linkID").Return(models.Link{
		DestinationURL: "https://test-url.example.com",
		Parameters: models.Parameters{
			"param1": []string{"value1"},
			"param2": []string{"value2", "value3"},
		},
		AllowParametersOverride: true,
	}, nil).Once()

	s.statsRepoM.On("LogRequest", models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
		Parameters: models.Parameters{
			"reqParam1": []string{"reqValue1"},
			"reqParam2": []string{"reqValue2", "reqValue3"},
		},
	}).Return(nil).Once()

	link, err := s.svc.Redirect(s.ctx, models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
		Parameters: models.Parameters{
			"reqParam1": []string{"reqValue1"},
			"reqParam2": []string{"reqValue2", "reqValue3"},
		},
	})
	s.Require().NoError(err)
	s.Require().Equal(
		"https://test-url.example.com?reqParam1=reqValue1&reqParam2=reqValue2&reqParam2=reqValue3",
		link,
	)
}

func (s *serviceTestSuite) TestRedirectStatsRepoErrorNoImpact() {
	s.mdRepoM.On("ResolveLink", "linkID").Return(models.Link{
		DestinationURL: "https://test-url.example.com",
		Parameters: models.Parameters{
			"param1": []string{"value1"},
			"param2": []string{"value2", "value3"},
		},
		AllowParametersOverride: true,
	}, nil).Once()

	s.statsRepoM.On("LogRequest", models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
		Parameters: models.Parameters{
			"reqParam1": []string{"reqValue1"},
			"reqParam2": []string{"reqValue2", "reqValue3"},
		},
	}).Return(errors.New("blah")).Once()

	link, err := s.svc.Redirect(s.ctx, models.Request{
		Timestamp: time.Date(2024, 11, 9, 1, 2, 3, 4, time.UTC),
		LinkID:    "linkID",
		ClientIP:  "127.3.2.1",
		CookieID:  "test-cookie-id",
		UserAgent: "some-app/1.0",
		Parameters: models.Parameters{
			"reqParam1": []string{"reqValue1"},
			"reqParam2": []string{"reqValue2", "reqValue3"},
		},
	})
	s.Require().NoError(err)
	s.Require().Equal(
		"https://test-url.example.com?reqParam1=reqValue1&reqParam2=reqValue2&reqParam2=reqValue3",
		link,
	)
}

// Definitions ...

type serviceTestSuite struct {
	suite.Suite

	ctx        context.Context
	mdRepoM    *mdRepoMock.Mock
	statsRepoM *statsRepoMock.Mock
	svc        Service
}

func (s *serviceTestSuite) SetupTest() {
	s.ctx = context.TODO()
	s.mdRepoM = mdRepoMock.New()
	s.statsRepoM = statsRepoMock.New()
	s.svc = New(s.mdRepoM, s.statsRepoM)
}

func (s *serviceTestSuite) TearDownTest() {
	s.ctx = nil
	s.mdRepoM.AssertExpectations(s.T())
	s.svc = nil
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, &serviceTestSuite{})
}
