package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mainflux/mainflux/logger"
	"github.com/mteodor/edgex-app/exapp"
	"github.com/mteodor/edgex-app/exapp/postgres"
	nats "github.com/nats-io/go-nats"
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
	// Connect to a NATS server

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	logger.Info(fmt.Sprintf("connecting %s", defNatsURL))
	nc, err := nats.Connect(defNatsURL)
	if err != nil {
		logger.Error("Failed to connect")
	}
	defer closeConn(nc)
	// Simple Async Subscriber
	nc.Subscribe("out.unknown", exapp.NatsMSGHandler())

	errs := make(chan error, 2)

	go func() {
		logger.Info(fmt.Sprintf("edgex-app started, exposed port %s", cfg.Port))
		errs <- exapp.InitHTTP(cfg.Port)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("edgex-app terminated: %s", err))

}

func closeConn(nc *nats.Conn) {
	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	fmt.Printf("closing down")
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

func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sql.DB {
	db, err := postgres.Connect(dbConfig)

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to postgres: %s", err))
		os.Exit(1)
	}
	return db
}

//geting enviroment
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
