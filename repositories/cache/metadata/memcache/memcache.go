package memcache

import (
	"context"
	"time"

	memcacheCli "github.com/bradfitz/gomemcache/memcache"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/teran/linker/models"
	protov1 "github.com/teran/linker/repositories/cache/metadata/memcache/proto/v1"
	"github.com/teran/linker/repositories/metadata"
)

var _ metadata.Repository = (*memcache)(nil)

type memcache struct {
	cli  *memcacheCli.Client
	repo metadata.Repository
	ttl  time.Duration
}

func New(cli *memcacheCli.Client, repo metadata.Repository, ttl time.Duration) metadata.Repository {
	return &memcache{
		cli:  cli,
		repo: repo,
		ttl:  ttl,
	}
}

func (m *memcache) ResolveLink(ctx context.Context, linkID string) (models.Link, error) {
	cacheKey := linkID
	item, err := m.cli.Get(cacheKey)
	if err != nil {
		if errors.Is(err, memcacheCli.ErrCacheMiss) {
			log.WithFields(log.Fields{
				"key": cacheKey,
			}).Tracef("cache miss")

			v, err := m.repo.ResolveLink(ctx, linkID)
			if err != nil {
				return models.Link{}, err
			}

			pv, err := protov1.NewLink(v)
			if err != nil {
				return models.Link{}, errors.Wrap(err, "error creating DTO")
			}

			payload, err := proto.Marshal(pv)
			if err != nil {
				return models.Link{}, errors.Wrap(err, "error marshaling DTO")
			}

			if err := m.cli.Set(&memcacheCli.Item{
				Key:        cacheKey,
				Expiration: int32(m.ttl.Seconds()),
				Value:      payload,
			}); err != nil {
				return models.Link{}, errors.Wrap(err, "error storing DTO")
			}

			return v, err
		}
		return models.Link{}, err
	}

	log.WithFields(log.Fields{
		"key": cacheKey,
	}).Tracef("cache hit")

	rv := &protov1.Link{}
	err = proto.Unmarshal(item.Value, rv)
	if err != nil {
		return models.Link{}, err
	}

	return rv.ToSvc()
}
