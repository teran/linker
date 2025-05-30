package migrations

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	postgresApp "github.com/teran/go-docker-testsuite/applications/postgres"
)

const postgreSQLImage = "index.docker.io/library/postgres:17"

func init() {
	log.SetLevel(log.TraceLevel)
}

func (s *migrateTestSuite) TestMigrateTwice() {
	err := migrateUpWithMigrationsPath(s.postgresDBApp.MustDSN("test_db"), "file://sql")
	s.Require().NoError(err)

	err = migrateUpWithMigrationsPath(s.postgresDBApp.MustDSN("test_db"), "file://sql")
	s.Require().NoError(err)

	err = migrateDownWithMigrationsPath(s.postgresDBApp.MustDSN("test_db"), "file://sql")
	s.Require().NoError(err)
}

// Definitions ...
type migrateTestSuite struct {
	suite.Suite

	ctx           context.Context
	postgresDBApp postgresApp.PostgreSQL
}

func (s *migrateTestSuite) SetupTest() {
	s.ctx = context.Background()

	app, err := postgresApp.NewWithImage(s.ctx, postgreSQLImage)
	s.Require().NoError(err)

	s.postgresDBApp = app

	err = s.postgresDBApp.CreateDB(s.ctx, "test_db")
	s.Require().NoError(err)
}

func (s *migrateTestSuite) TearDownTest() {
	err := s.postgresDBApp.Close(s.ctx)
	s.Require().NoError(err)
}

func TestMigrateTestSuite(t *testing.T) {
	suite.Run(t, &migrateTestSuite{})
}
