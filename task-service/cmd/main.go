package main

import (
	"log"
	"net"
	"time"

	pb "github.com/sejamuchhal/taskhub/task-service/proto"
	"github.com/sejamuchhal/taskhub/task-service/server"
	"github.com/sejamuchhal/taskhub/task-service/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	listener, err := net.Listen("tcp", config.GRPCAddress)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", config.GRPCAddress, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 3 * time.Minute,
		Timeout:           10 * time.Second,
		MaxConnectionAge:  50 * time.Minute,
		Time:              10 * time.Minute,
	}))

	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	pb.RegisterTaskServiceServer(grpcServer, srv)

	log.Printf("Starting gRPC server on %s...", config.GRPCAddress)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
