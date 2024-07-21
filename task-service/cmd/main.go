package main

import (
	"log"
	"net"
	"time"

	"github.com/sejamuchhal/taskhub/task-service/common"
	pb "github.com/sejamuchhal/taskhub/task-service/pb/task"
	"github.com/sejamuchhal/taskhub/task-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Print("Starting task server")
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", "0.0.0.0:8080", err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 3 * time.Minute,
		Timeout:           10 * time.Second,
		MaxConnectionAge:  50 * time.Minute,
		Time:              10 * time.Minute,
	}))

	reflection.Register(grpcServer)
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
