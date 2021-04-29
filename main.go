package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/proemergotech/log/v3"
	"github.com/proemergotech/log/v3/zaplog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/artofimagination/mysql-resources-db-go-service/config"
	"github.com/artofimagination/mysql-resources-db-go-service/initialization"
)

type contextMapper struct{}

func (cl contextMapper) Values(_ context.Context) map[string]string {
	// nil mapper
	// todo: could add fields to context and log them when necessary just by extract and adding them here to the logger
	return map[string]string{}
}

func ContextMapper() log.ContextMapper {
	return contextMapper{}
}

func main() {

	if err := zap.RegisterEncoder(
		zaplog.EncoderType,
		zaplog.NewEncoder([]string{
			log.AppName,
			log.AppVersion,
		}),
	); err != nil {
		panic(fmt.Sprintf("Couldn't create logger, error: %v", err))
	}

	zapConf := zap.NewProductionConfig()
	zapConf.Encoding = zaplog.EncoderType

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	zapLogLevel := new(zapcore.Level)
	if err := zapLogLevel.Set(logLevel); err != nil {
		panic(fmt.Sprintf("Invalid log level: %s", logLevel))
	}
	zapConf.Level = zap.NewAtomicLevelAt(*zapLogLevel)

	zapLog, err := zapConf.Build()
	if err != nil {
		panic(fmt.Sprintf("Couldn't create logger, error: %v", err))
	}
	zapLog = zapLog.With(
		zap.String(log.AppName, config.AppName),
		zap.String(log.AppVersion, config.AppVersion),
	)

	log.SetGlobalLogger(zaplog.NewLogger(zapLog, ContextMapper()))

	defer func() {
		if err := recover(); err != nil {
			log.Error(context.Background(), "Service panicked", "error", errors.Errorf("%+v", err))
			os.Exit(1)
		}
	}()

	initialization.Execute()
}
