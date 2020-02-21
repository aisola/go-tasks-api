package main

import (
	"log"
	"net/http"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"example.com/tasks"
	"example.com/tasks/sqlite"
	"example.com/tasks/taskhttp"
)

func init() {
	pflag.StringP("bind", "b", ":5000", "The interface and port on which to serve.")
	pflag.StringP("database", "d", ":memory:", "The path to the sqlite3 database.")

	viper.BindPFlag("bind", pflag.Lookup("bind"))
	viper.BindPFlag("database", pflag.Lookup("database"))
}

func initializeLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err)
	}

	return logger
}

func initializeRepository(logger *zap.Logger, database string) tasks.TaskRepository {
	repo, err := sqlite.New(database)
	if err != nil {
		logger.Error("failed to initialize database", zap.Error(err))
		os.Exit(1)
	}

	return repo
}

func main() {
	pflag.Parse()

	logger := initializeLogger()
	defer logger.Sync()
	repo := initializeRepository(logger, viper.GetString("database"))

	handler := taskhttp.New(logger.Named("tasks"), repo)

	logger.Info("I'm Listening", zap.String("bind", viper.GetString("bind")))
	if err := http.ListenAndServe(viper.GetString("bind"), handler); err != nil {
		logger.Error("failed to listen and serve",
			zap.String("bind", viper.GetString("bind")),
			zap.Error(err),
		)
		os.Exit(1)
	}
}
