package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/mainflux/mainflux/logger"
	"github.com/mteodor/edgex-app/exapp"
	"github.com/mteodor/edgex-app/exapp/api"
	httpapi "github.com/mteodor/edgex-app/exapp/api/http"
	"github.com/mteodor/edgex-app/exapp/postgres"
	nats "github.com/nats-io/go-nats"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	defNatsURL       = nats.DefaultURL
	defLogLevel      = "error"
	defPort          = "8000"
	envNatsURL       = "MF_NATS_URL"
	envLogLevel      = "MF_EDGEX_APP_LOG_LEVEL"
	envPort          = "MF_EDGEX_APP_PORT"
	defDBHost        = "localhost"
	defDBPort        = "5432"
	defDBUser        = "mainflux"
	defDBPass        = "mainflux"
	defDBName        = "events"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	envDBHost        = "MF_EDGEX_DB_HOST"
	envDBPort        = "MF_EDGEX_DB_PORT"
	envDBUser        = "MF_EDGEX_DB_USER"
	envDBPass        = "MF_EDGEX_DB_PASS"
	envDBName        = "MF_EDGEX_DB"
	envDBSSLMode     = "MF_EDGEX_DB_SSL_MODE"
	envDBSSLCert     = "MF_EDGEX_DB_SSL_CERT"
	envDBSSLKey      = "MF_EDGEX_DB_SSL_KEY"
	envDBSSLRootCert = "MF_EDGEX_DB_SSL_ROOT_CERT"
	envMqttHost      = "MF_EDGEX_MQTT_HOST"
	envMqttPassword  = "MF_EDGEX_MQTT_PASS"
	envMqttUsername  = "MF_EDGEX_MQTT_USERNAME"
	envMqttSndTopic  = "MF_EDGEX_MQTT_TOPIC"
	envMqttClientID  = "MF_EDGEX_MQTT_CLIENT_ID"
	topicUnknown     = "out.unknown"
)

type config struct {
	dbConfig postgres.Config
	NatsURL  string
	LogLevel string
	Port     string
}

// Connect to a server

func main() {

	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// Connect to a DB server
	db := connectToDB(cfg.dbConfig, logger)
	if db == nil {
		log.Fatalf("cannot connect to db")
	}
	defer db.Close()

	// Make a MQTT publish connection
	mfMqttHost := Env(envMqttHost, "")
	mfMqttPassword := Env(envMqttPassword, "")
	mfMqttUsername := Env(envMqttUsername, "")
	mfMqttClientID := Env(envMqttClientID, "")
	mfMqttSndTopic := Env(envMqttSndTopic, "")
	if mfMqttPassword != "" && mfMqttUsername != "" && mfMqttSndTopic != "" {
		_, err := api.NewMQTTPublisher(&logger, mfMqttHost, mfMqttUsername, mfMqttPassword, mfMqttClientID, mfMqttSndTopic)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create mqtt publisher: %s", err))
		}
	}

	// create a service
	svc := newService(db, logger)

	// connect to NATS
	connectToNats(svc, logger)

	// configure endpoints and start http server
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), httpapi.MakeHandler(svc, logger))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to init http server on port %s: ", cfg.Port))
		return
	}

	logger.Info("init server done")
	errs := make(chan error, 2)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("edgex-app terminated: %s", err))

}

func closeConn(nc *nats.Conn, logger logger.Logger) {
	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	logger.Info("closing down")
	if nc == nil {
		return
	}

	nc.Drain()
	nc.Close()
}

func loadConfig() config {

	dbConfig := postgres.Config{
		Host:        Env(envDBHost, defDBHost),
		Port:        Env(envDBPort, defDBPort),
		User:        Env(envDBUser, defDBUser),
		Pass:        Env(envDBPass, defDBPass),
		Name:        Env(envDBName, defDBName),
		SSLMode:     Env(envDBSSLMode, defDBSSLMode),
		SSLCert:     Env(envDBSSLCert, defDBSSLCert),
		SSLKey:      Env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: Env(envDBSSLRootCert, defDBSSLRootCert),
	}

	return config{
		NatsURL:  Env(envNatsURL, defNatsURL),
		LogLevel: Env(envLogLevel, defLogLevel),
		Port:     Env(envPort, defPort),
		dbConfig: dbConfig,
	}
}
func connectToNats(svc exapp.Service, logger logger.Logger) {

	logger.Info(fmt.Sprintf("connecting %s", defNatsURL))

	nc, err := nats.Connect(defNatsURL)
	defer closeConn(nc, logger)

	if err != nil {
		logger.Error("Failed to connect to nats")
	} else {
		nc.Subscribe(topicUnknown, exapp.NatsMSGHandler(svc))
	}
}
func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sql.DB {
	db, err := postgres.Connect(dbConfig, logger)

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to postgres: %s", err))
		os.Exit(1)
	}
	logger.Info("connected to database")
	return db
}

// Env geting enviroment
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func newService(db *sql.DB, logger logger.Logger) exapp.Service {

	eventsRepository := postgres.New(db, logger)
	svc := exapp.New(eventsRepository, logger)
	svc = api.LoggingMiddleware(svc)

	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "events",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "events",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}
