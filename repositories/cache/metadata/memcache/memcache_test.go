package memcache

import (
	"context"
	"testing"
	"time"

	memcacheCli "github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	memcacheApp "github.com/teran/go-docker-testsuite/applications/memcache"
	"github.com/teran/linker/models"
	"github.com/teran/linker/repositories/metadata"
	repoM "github.com/teran/linker/repositories/metadata/mock"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func (s *memcacheTestSuite) TestResolveLink() {
	const linkID = "test-link-id"
	s.repoMock.On("ResolveLink", linkID).Return(models.Link{
		DestinationURL: "https://test-url.example.com",
		Parameters: models.Parameters{
			"key1": []string{"value1"},
			"key2": []string{"value2"},
		},
		AllowParametersOverride: true,
	}, nil).Once()

	link, err := s.cache.ResolveLink(s.ctx, linkID)
	s.Require().NoError(err)
	s.Require().Equal(models.Link{
		DestinationURL: "https://test-url.example.com",
		Parameters: models.Parameters{
			"key1": []string{"value1"},
			"key2": []string{"value2"},
		},
		AllowParametersOverride: true,
	}, link)

	link, err = s.cache.ResolveLink(s.ctx, linkID)
	s.Require().NoError(err)
	s.Require().Equal(models.Link{
		DestinationURL: "https://test-url.example.com",
		Parameters: models.Parameters{
			"key1": []string{"value1"},
			"key2": []string{"value2"},
		},
		AllowParametersOverride: true,
	}, link)
}

// Definitions ...
type memcacheTestSuite struct {
	suite.Suite

	ctx      context.Context
	cache    metadata.Repository
	repoMock *repoM.Mock

	memcachedApp memcacheApp.Memcache
}

func (s *memcacheTestSuite) SetupSuite() {
	s.ctx = context.TODO()

	var err error
	s.memcachedApp, err = memcacheApp.New(s.ctx)
	s.Require().NoError(err)
}

func (s *memcacheTestSuite) SetupTest() {
	s.repoMock = repoM.New()

	url, err := s.memcachedApp.GetEndpointAddress()
	s.Require().NoError(err)

	cli := memcacheCli.New(url)
	s.cache = New(cli, s.repoMock, 3*time.Second)
}

func (s *memcacheTestSuite) TearDownTest() {
	s.repoMock.AssertExpectations(s.T())
}

func (s *memcacheTestSuite) TearDownSuite() {
	s.memcachedApp.Close(s.ctx)
}

func TestMemcacheTestSuite(t *testing.T) {
	suite.Run(t, &memcacheTestSuite{})
}
