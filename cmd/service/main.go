package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/andrii-lavreniuk/oapi-fiber-example/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	logger := newLogger(cfg.GetLog()).WithField("service", cfg.GetServiceName()).Logger

	app, cleanup, err := wireApp(cfg.GetServer(), cfg.GetData(), logger)
	if err != nil {
		logger.WithError(err).Fatal("wire build app")
	}

	defer cleanup()

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.GetServer().Port)); err != nil {
			logger.WithError(err).Fatal("failed to start server")
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received

	logger.Debug("gracefully shutting down...")

	_ = app.Shutdown()
}

func newLogger(c config.Log) *logrus.Logger {
	lg := logrus.New()
	lg.SetLevel(logrus.WarnLevel)
	lg.SetReportCaller(true)

	if c.Format == "text" {
		lg.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
				fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
				return "", fileName
			},
		})
	} else {
		lg.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
				fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
				return "", fileName
			},
		})
	}

	ll, err := logrus.ParseLevel(c.Level)
	if err != nil {
		lg.WithError(err).Warnf("failed to parse log level: %s", c.Level)
	} else {
		lg.SetLevel(ll)
	}

	return lg
}
