package service

import (
	"context"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/teran/linker/models"
	"github.com/teran/linker/repositories/metadata"
	"github.com/teran/linker/repositories/stats"
)

type Service interface {
	Redirect(ctx context.Context, req models.Request) (string, error)
}

type service struct {
	mdRepo    metadata.Repository
	statsRepo stats.Repository
}

func New(mdRepo metadata.Repository, statsRepo stats.Repository) Service {
	return &service{
		mdRepo:    mdRepo,
		statsRepo: statsRepo,
	}
}

func (s *service) Redirect(ctx context.Context, req models.Request) (string, error) {
	link, err := s.mdRepo.ResolveLink(ctx, req.LinkID)
	if err != nil {
		return "", mapRepositoryErrors(err)
	}

	if err := s.statsRepo.LogRequest(ctx, req); err != nil {
		log.Warnf("error storing request into stats repository: `%s`", err)
	}

	lnk, err := url.Parse(link.DestinationURL)
	if err != nil {
		return "", err
	}

	lnk.RawQuery = link.Parameters.String()

	if link.AllowParametersOverride && len(req.Parameters) > 0 {
		lnk.RawQuery = req.Parameters.String()
	}

	return lnk.String(), nil
}
