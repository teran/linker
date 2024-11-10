package postgresql

import (
	"github.com/teran/linker/models"
	"github.com/teran/linker/repositories/metadata"
)

func (s *postgreSQLRepositoryTestSuite) TestResolveLinkNotFound() {
	_, err := s.repo.ResolveLink(s.ctx, "testID")
	s.Require().Error(err)
	s.Require().Equal(metadata.ErrNotFound, err)
}

func (s *postgreSQLRepositoryTestSuite) TestResolveLink() {
	sql, args, err := psql.
		Insert("links").
		Columns(
			"link_id",
			"destination_url",
			"parameters",
			"allow_parameters_override",
		).
		Values(
			"testID",
			"https://test-url.example.com",
			models.Parameters{},
			true,
		).
		ToSql()
	s.Require().NoError(err)

	_, err = s.db.ExecContext(s.ctx, sql, args...)
	s.Require().NoError(err)

	l, err := s.repo.ResolveLink(s.ctx, "testID")
	s.Require().NoError(err)
	s.Require().Equal(models.Link{
		DestinationURL:          "https://test-url.example.com",
		Parameters:              models.Parameters{},
		AllowParametersOverride: true,
	}, l)
}
