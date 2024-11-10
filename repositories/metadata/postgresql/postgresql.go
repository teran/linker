package postgresql

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/teran/linker/repositories/metadata"
)

var (
	_ metadata.Repository = (*repository)(nil)

	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type repository struct {
	master  *sql.DB
	replica *sql.DB
}

func mapSQLErrors(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return metadata.ErrNotFound
	}

	return err
}

func New(master, replica *sql.DB) metadata.Repository {
	return &repository{
		master:  master,
		replica: replica,
	}
}

type queryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func selectQueryRow(ctx context.Context, db queryRunner, q sq.SelectBuilder) (sq.RowScanner, error) {
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	start := time.Now()
	defer func() {
		log.WithFields(log.Fields{
			"query":    sql,
			"args":     args,
			"duration": time.Now().Sub(start),
		}).Debug("SQL query executed")
	}()

	return db.QueryRowContext(ctx, sql, args...), nil
}
