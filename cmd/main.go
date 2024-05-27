package main

import (
	"context"
	"fmt"
	"monitoring-system/server/cmd/factory"
	"monitoring-system/server/cmd/server"
	"monitoring-system/server/config"
	"monitoring-system/server/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// type Application struct {
// 	logger  logger.Logger
// 	storage storage.Storage
// 	config  *config.Config
// 	ctx     context.Context
// 	cm      camera_manager.CameraManager
// 	sqlDB   *sql.DB
// 	modules *modules.Modules
// }

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

	awsConfig, err := config.LoadAwsConfig(ctx, appConfig.Aws, logger)
	if err != nil {
		logger.Error("Error loading AWS configuration %v", err)
		return
	}

	factory, err := factory.New(ctx, *awsConfig, *appConfig)
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
