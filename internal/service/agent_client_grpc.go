package service

import (
	"fmt"
	pb "metrics/internal/proto/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	pb.MetricsClient
	conn *grpc.ClientConn
}

func (dc *GrpcClient) Close() error {
	err := dc.conn.Close()
	if err != nil {
		return fmt.Errorf("close grpc client: %w", err)
	}
	return nil
}

func NewGrpcClient() (*GrpcClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("localhost:8081", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new client: %w", err)
	}

	client := pb.NewMetricsClient(conn)

	return &GrpcClient{
		conn:          conn,
		MetricsClient: client,
	}, nil
}
