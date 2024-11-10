package postgresql

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	postgresApp "github.com/teran/go-docker-testsuite/applications/postgres"

	"github.com/teran/linker/repositories/metadata"
	"github.com/teran/linker/repositories/metadata/postgresql/migrations"
)

const postgreSQLImage = "index.docker.io/library/postgres:17"

// Definitions ...
type postgreSQLRepositoryTestSuite struct {
	suite.Suite

	ctx           context.Context
	postgresDBApp postgresApp.PostgreSQL
	db            *sql.DB
	repo          metadata.Repository
}

func (s *postgreSQLRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	app, err := postgresApp.NewWithImage(s.ctx, postgreSQLImage)
	s.Require().NoError(err)

	s.postgresDBApp = app
}

func (s *postgreSQLRepositoryTestSuite) SetupTest() {
	err := s.postgresDBApp.CreateDB(s.ctx, "test_db")
	s.Require().NoError(err)

	dsn, err := s.postgresDBApp.DSN("test_db")
	s.Require().NoError(err)

	err = migrations.MigrateUp(dsn)
	s.Require().NoError(err)

	db, err := sql.Open("postgres", dsn)
	s.Require().NoError(err)

	s.db = db

	s.repo = New(s.db, s.db)
}

func (s *postgreSQLRepositoryTestSuite) TearDownTest() {
	s.repo = nil

	_, err := s.db.ExecContext(
		s.ctx,
		`SELECT
			pg_terminate_backend(pg_stat_activity.pid)
		FROM
			pg_stat_activity
		WHERE
			pg_stat_activity.datname = 'test_db'
		AND
			pid != pg_backend_pid();
	`)
	s.Require().NoError(err)

	err = s.db.Close()
	s.Require().NoError(err)

	err = s.postgresDBApp.DropDB(s.ctx, "test_db")
	s.Require().NoError(err)
}

func (s *postgreSQLRepositoryTestSuite) TearDownSuite() {
	err := s.postgresDBApp.Close(s.ctx)
	s.Require().NoError(err)
}

func TestPostgreSQLRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &postgreSQLRepositoryTestSuite{})
}
