package db

import (
	"context"
	"os"
	
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	log "github.com/sirupsen/logrus"
)

var DBpool *pgxpool.Pool

func Start() error {
	// Open database log file
	logFile, err := os.OpenFile("db.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	config, err := pgxpool.ParseConfig("postgres://localhost:5432/nail")
	if err != nil {
		return err
	}
	pgxLogger := &log.Logger{
		Out:          logFile,
		Formatter:    new(log.TextFormatter),
		Hooks:        make(log.LevelHooks),
		Level:        log.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
    }
	config.ConnConfig.Logger = logrusadapter.NewLogger(pgxLogger)
	config.ConnConfig.LogLevel = pgx.LogLevelInfo

	DBpool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return err
	}
	return nil
}
