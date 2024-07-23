package main

import (
	"log"

	"github.com/sejamuchhal/taskhub/gateway/common"
	srv "github.com/sejamuchhal/taskhub/gateway/server"
)

func main() {
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	log.Println("Starting HTTP server")

	server, err := srv.NewServer(config)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Cannot start server: %v", err)
	}
}
