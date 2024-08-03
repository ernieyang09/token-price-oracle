package main

import (
	"oracle-go/pkg/env"
	"oracle-go/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
)


func main() {
	logPath := env.GetConfigValue("LOG_PATH", "price_syncer.log")
	logger.InitLogger(logPath)
	logger.Info("Starting cron job")
	c := cron.New(cron.WithSeconds())
    
    SetupJobs(c)

    c.Start()

	signalChannel := make(chan os.Signal, 1)

	// Notify the signal channel on SIGINT (Ctrl+C) and SIGTERM
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal
	<-signalChannel

	// Log termination message
	logger.Info("Received termination signal. Shutting down...")

	// Stop the cron scheduler
	c.Stop()

	// Optionally, you can wait for a short period to ensure logs are flushed
	logger.Info("Shutdown complete")
	
}