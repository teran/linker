package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo-contrib/echoprometheus"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/teran/appmetrics"
	"github.com/teran/go-random"
	"golang.org/x/sync/errgroup"

	"github.com/teran/linker/presenter"
	mdRepoPostgres "github.com/teran/linker/repositories/metadata/postgresql"
	statsRepoKafka "github.com/teran/linker/repositories/stats/kafka"
	"github.com/teran/linker/service"
)

var (
	appVersion     = "n/a (dev build)"
	buildTimestamp = "undefined"
)

type config struct {
	Addr        string `envconfig:"ADDR" default:":8080"`
	MetricsAddr string `envconfig:"METRICS_ADDR" default:":8081"`

	KafkaBrokers    []string `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaStatsTopic string   `envconfig:"KAFKA_STATS_TOPIC" required:"true"`

	MetadataDBMasterDSN  string `envconfig:"METADATADB_MASTER_DSN" required:"true"`
	MetadataDBReplicaDSN string `envconfig:"METADATADB_REPLICA_DSN" required:"true"`

	LogLevel log.Level `envconfig:"LOG_LEVEL" default:"info"`

	Domain     string `envconfig:"SERVICE_DOMAIN" required:"true"`
	CookieName string `envconfig:"SERVICE_COOKIE_NAME" required:"true"`
}

func main() {
	var cfg config
	envconfig.MustProcess("", &cfg)

	log.SetLevel(cfg.LogLevel)
	sarama.Logger = log.StandardLogger()

	lf := new(log.TextFormatter)
	lf.FullTimestamp = true
	log.SetFormatter(lf)

	log.Infof("Initializing linker-redirector (%s @ %s) ...", appVersion, buildTimestamp)

	g, _ := errgroup.WithContext(context.Background())

	dbMaster, err := sql.Open("postgres", cfg.MetadataDBMasterDSN)
	if err != nil {
		panic(err)
	}

	dbReplica, err := sql.Open("postgres", cfg.MetadataDBReplicaDSN)
	if err != nil {
		panic(err)
	}

	mdRepo := mdRepoPostgres.New(dbMaster, dbReplica)

	kcfg := newProducerConfig()
	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, kcfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Warnf("Failed to close Kafka producer cleanly: %s", err)
		}
	}()

	statsRepo := statsRepoKafka.New(producer, cfg.KafkaStatsTopic)

	svc := service.New(mdRepo, statsRepo)
	handlers := presenter.New(&presenter.Config{
		Domain:     cfg.Domain,
		CookieName: cfg.CookieName,
		URLService: svc,
		CookieIDGenerator: func() string {
			return random.String(append(random.AlphaNumeric, '.', '_', '-', '/', ','), 31)
		},
		TimeNowProvider: func() time.Time {
			return time.Now()
		},
	})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(echoprometheus.NewMiddleware("showcase"))
	e.Use(middleware.Recover())
	handlers.Register(e)

	g.Go(func() error {
		srv := &http.Server{
			Addr:    cfg.Addr,
			Handler: e,
		}

		return srv.ListenAndServe()
	})

	me := echo.New()
	if cfg.LogLevel >= log.DebugLevel {
		me.Use(middleware.Logger())
	}
	me.Use(echoprometheus.NewMiddleware("showcase_metrics"))
	me.Use(middleware.Recover())

	checkFn := func() error {
		if err := dbMaster.Ping(); err != nil {
			return err
		}

		if err := dbReplica.Ping(); err != nil {
			return err
		}

		return nil
	}

	metrics := appmetrics.New(checkFn, checkFn, checkFn)
	metrics.Register(me)

	g.Go(func() error {
		srv := http.Server{
			Addr:    cfg.MetricsAddr,
			Handler: me,
		}

		return srv.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		panic(err)
	}
}

func newProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	return config
}
