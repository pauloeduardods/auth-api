package main

import (
	"context"
	"database/sql"
	"fmt"
	"monitoring-system/server/src/cmd/factory"
	"monitoring-system/server/src/cmd/server"
	"monitoring-system/server/src/config"
	"monitoring-system/server/src/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
)

func main() {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("Error loading configuration %v", err)
		return
	}

	logger, err := logger.NewLogger(appConfig.Env)
	if err != nil {
		fmt.Printf("Error creating logger %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logger.Info("Received signal: %v", sig)
		cancel()
	}()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", appConfig.Sql.Host, appConfig.Sql.Port, appConfig.Sql.User, appConfig.Sql.Password, appConfig.Sql.Database)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Error connecting to database %v", err)
		return
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		logger.Error("Error pinging database %v", err)
		return
	}

	awsConfig, err := config.LoadAwsConfig(ctx, appConfig.Aws, logger)
	if err != nil {
		logger.Error("Error loading AWS configuration %v", err)
		return
	}

	factory, err := factory.New(ctx, logger, *awsConfig, *appConfig, db)
	if err != nil {
		logger.Error("Error creating factory %v", err)
		return
	}

	server := server.New(ctx, awsConfig, appConfig, logger, factory)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Start(); err != nil {
			logger.Error("Error starting server %v", err)
		}
	}()

	wg.Wait()

	<-ctx.Done()
	os.Exit(0)
}
