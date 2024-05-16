package main

import (
	"context"
	"fmt"
	"monitoring-system/server/cmd/server"
	"monitoring-system/server/config"
	"monitoring-system/server/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	os.Exit(start())
}

func start() int {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return 1
	}

	logger, err := logger.NewLogger(appConfig.AppEnv)
	if err != nil {
		fmt.Println("Error setting up the logger:", err)
		return 1
	}
	//TODO REMOVE
	logger.GetZapLogger().Info("Starting server", zap.String("host", appConfig.Host), zap.Int("port", appConfig.Port))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logger.Info("Received signal: %v", sig)
		cancel()
	}()

	awsConfig, err := config.NewAWSConfig(ctx, appConfig, logger)
	if err != nil {
		logger.Info("Error creating AWS config %v", err)
		return 1
	}

	s := server.New(ctx, awsConfig, appConfig, logger)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			logger.Info("Error starting server %v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		if err := s.Stop(); err != nil {
			logger.Info("Error stopping server %v", err)
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}
