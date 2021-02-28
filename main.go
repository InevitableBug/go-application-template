package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"time"

	"bug.geek.nz/go-application-template/config"
	"bug.geek.nz/go-application-template/server"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	initialiseLog()

	srv := server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// block until signal received
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_ = srv.Shutdown(ctx)

	log.Info("shutting down")
	os.Exit(0)

}

func initialiseLog() {

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	fileWriter := &lumberjack.Logger{
		Filename: config.Instance.Log.LogFile,
		MaxSize:  20, // megabytes
	}

	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	log.SetOutput(multiWriter)

	switch config.Instance.Log.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}
