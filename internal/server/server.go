package server

import (
	"fmt"
	"metrics/internal/config"
	"metrics/internal/repository"
	"metrics/internal/router"
	"metrics/internal/service"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "net/http/pprof"

	"go.uber.org/zap"

	"metrics/internal/handlers/rpc"
	pb "metrics/internal/proto/v1"
)

func ConfigureServerHandler(
	memStorage repository.MetricStorage,
	cfg *config.ServerConfig,
	logger *zap.SugaredLogger,
) (*http.Server, error) {
	handlerLogger := logger.With("r", "r")

	r := router.ConfigureServerHandler(memStorage, cfg, logger)
	handlerLogger.Infow(
		"Starting server",
		"addr", cfg.Address,
	)
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		return srv, fmt.Errorf("listen and server has failed: %w", err)
	}

	return srv, nil
}

func InitPprof() (*http.Server, error) {
	pprofServer := &http.Server{
		Addr: "0.0.0.0:6060",
	}
	if err := pprofServer.ListenAndServe(); err != nil {
		return nil, fmt.Errorf("listen and server has failed: %w", err)
	}

	return pprofServer, nil
}

func Serve(memStorage repository.MetricStorage, logger *zap.SugaredLogger) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		return nil, fmt.Errorf("failed to run gRPC server: %w", err)
	}
	logger.Infow(
		"Starting grpc server",
		"addr", "localhost:8081",
	)

	metricService := service.NewMetricService(memStorage, logger)

	grpcServer := grpc.NewServer()
	pb.RegisterMetricsServer(grpcServer, rpc.NewServer(metricService, logger))

	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		return nil, fmt.Errorf("failed to run gRPC server: %w", err)
	}

	return grpcServer, nil
}
