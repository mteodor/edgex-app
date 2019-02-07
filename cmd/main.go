package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/logger"
	"github.com/mteodor/edgex-app/exapp"
	nats "github.com/nats-io/go-nats"
)

const (
	defNatsURL  string = nats.DefaultURL
	defLogLevel string = "error"
	defPort     string = "8000"
	envNatsURL  string = "MF_NATS_URL"
	envLogLevel string = "MF_EDGEX_APP_LOG_LEVEL"
	envPort     string = "MF_EDGEX_APP_PORT"
)

type config struct {
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
	// Connect to a server

	logger.Info(fmt.Sprintf("connecting %s\n", defNatsURL))
	nc, err := nats.Connect(defNatsURL)
	if err != nil {
		logger.Error("Failed to connect\n")
	}
	defer closeConn(nc)
	// Simple Async Subscriber
	nc.Subscribe("out.unknown", exapp.MakeHandler())

	errs := make(chan error, 2)

	go func() {
		p := fmt.Sprintf(":%s", cfg.Port)
		logger.Info(fmt.Sprintf("edgex-app started, exposed port %s", cfg.Port))
		errs <- http.ListenAndServe(p, exapp.MakeHTTPHandler())
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
	return config{
		NatsURL:  mainflux.Env(envNatsURL, defNatsURL),
		LogLevel: mainflux.Env(envLogLevel, defLogLevel),
		Port:     mainflux.Env(envPort, defPort),
	}
}
