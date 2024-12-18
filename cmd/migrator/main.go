package main

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/teran/linker/repositories/metadata/postgresql/migrations"
)

var (
	appVersion     = "n/a (dev build)"
	buildTimestamp = "undefined"
)

type config struct {
	LogLevel log.Level `envconfig:"LOG_LEVEL" default:"info"`

	MetadataDBMasterDSN string `envconfig:"METADATADB_MASTER_DSN" required:"true"`
}

func main() {
	var cfg config
	envconfig.MustProcess("", &cfg)

	log.SetLevel(cfg.LogLevel)

	lf := new(log.TextFormatter)
	lf.FullTimestamp = true
	log.SetFormatter(lf)

	log.Infof("Initializing linker-migrator (%s @ %s) ...", appVersion, buildTimestamp)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return migrations.MigrateUp(cfg.MetadataDBMasterDSN)
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
