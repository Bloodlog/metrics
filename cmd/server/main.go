package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	config "metrics/internal/config/server"
	"metrics/internal/logger"
	"metrics/internal/repository"
	"metrics/internal/server"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"
)

var (
	version     = "N/A"
	buildTime   = "N/A"
	buildCommit = "N/A"
)

const (
	timeoutServerShutdown = time.Second * 5
	timeoutShutdown       = time.Second * 10
)

func main() {
	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", buildTime)
	fmt.Printf("Build commit: %s\n", buildCommit)
	loggerZap, err := logger.InitilazerLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	if err = run(loggerZap); err != nil {
		loggerZap.Fatal("Application failed", zap.Error(err))
	}
}

// @title Metrics API
// @version 1.0
// @description API для управления метриками
// @host 127.0.0.1:8080
// @BasePath /.
func run(loggerZap *zap.SugaredLogger) error {
	rootCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelCtx()

	g, ctx := errgroup.WithContext(rootCtx)

	context.AfterFunc(ctx, func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()

		<-ctx.Done()
		log.Fatal("failed to gracefully shutdown the service")
	})

	cfg, err := config.ParseFlags()
	if err != nil {
		loggerZap.Info(err.Error(), "failed to parse flags")
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	memStorage, err := repository.NewMetricStorage(ctx, cfg, loggerZap)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	g.Go(func() error {
		defer log.Print("closed DB")

		<-ctx.Done()

		memStorage.Shutdown(ctx)
		return nil
	})

	var httpServer *http.Server
	var pprofServer *http.Server
	var grpcServer *grpc.Server

	if cfg.Debug {
		g.Go(func() (err error) {
			pprofServer, err = server.InitPprof()
			if err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					return
				}
				return fmt.Errorf("listen and server has failed: %w", err)
			}
			return nil
		})
		g.Go(func() error {
			defer log.Print("PProf server has been shutdown")
			<-ctx.Done()

			shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
			defer cancelShutdownTimeoutCtx()

			if pprofServer != nil {
				if err := pprofServer.Shutdown(shutdownTimeoutCtx); err != nil {
					loggerZap.Info("PProf server gracefully stopped: %v", err)
				}
			}
			return nil
		})
	}

	g.Go(func() (err error) {
		httpServer, err = server.ConfigureServerHandler(memStorage, cfg, loggerZap)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			return fmt.Errorf("listen and server has failed: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		defer log.Print("server has been shutdown")
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()

		if httpServer != nil {
			if err := httpServer.Shutdown(shutdownTimeoutCtx); err != nil {
				loggerZap.Info("HTTP server Shutdown: %v", err)
			}
		}

		return nil
	})

	g.Go(func() (err error) {
		grpcServer, err = server.Serve(memStorage, loggerZap)
		if err != nil {
			return fmt.Errorf("listen and server grpc has failed: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		defer log.Print("grpc server has been shutdown")
		<-ctx.Done()

		if grpcServer != nil {
			grpcServer.GracefulStop()
		}

		log.Print("gRPC server has been shutdown")
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to wait for errgroup: %w", err)
	}

	return nil
}
