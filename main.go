package main

import (
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fmatoss/trumail/api"
	"github.com/fmatoss/trumail/config"
)

func main() {
	http.DefaultClient = &http.Client{
		Timeout: time.Duration(config.HTTPClientTimeout) * time.Second,
	}

	logger := logrus.New() // New Logger

	if strings.Contains(config.Env, "prod") {
		logger.Formatter = new(logrus.JSONFormatter)
	}
	l := logger.WithField("port", config.Port)

	r, s := api.Initialize(logger)

	l.Info("Binding all Trumail endpoints to the router")
	api.RegisterEndpoints(r, s)

	if config.ServeWeb {
		// Set all remaining paths to point to static files (must come after)
		r.HandleStatic("./web")
	}

	handleSignals()

	// Listen and Serve
	l.Info("Listening and Serving")
	r.ListenAndServe(config.Port)
}

func handleSignals() {
	sigChan := make(chan os.Signal, 3)
	go func() {
		for sig := range sigChan {
			if sig == syscall.SIGUSR1 {
				pprof.Lookup("goroutine").WriteTo(os.Stdout, 2)
			}
		}
	}()
	signal.Notify(sigChan, syscall.SIGUSR1)
}
