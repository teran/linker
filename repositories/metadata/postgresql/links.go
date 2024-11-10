package postgresql

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/teran/linker/models"
)

func (r *repository) ResolveLink(ctx context.Context, linkID string) (models.Link, error) {
	row, err := selectQueryRow(ctx, r.replica, psql.
		Select(
			"destination_url",
			"parameters",
			"allow_parameters_override",
		).
		From("links").
		Where(sq.Eq{
			"link_id": linkID,
		}),
	)
	if err != nil {
		return models.Link{}, mapSQLErrors(err)
	}

	l := models.Link{}
	if err := row.Scan(
		&l.DestinationURL,
		&l.Parameters,
		&l.AllowParametersOverride,
	); err != nil {
		return models.Link{}, mapSQLErrors(err)
	}

	return l, nil
}
