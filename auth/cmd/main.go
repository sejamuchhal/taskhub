package main

import (
	"log"
	"net"
	"time"

	"github.com/sejamuchhal/taskhub/auth/common"
	pb "github.com/sejamuchhal/taskhub/auth/pb"
	"github.com/sejamuchhal/taskhub/auth/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Print("Starting task server")
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

	reflection.Register(grpcServer)
	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	pb.RegisterAuthServiceServer(grpcServer, srv)

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	log.Println("Starting gRPC server on port 4040...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
