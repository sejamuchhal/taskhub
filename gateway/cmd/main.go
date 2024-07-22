package main

import (
	"fmt"
	"log"

	"github.com/sejamuchhal/taskhub/gateway/common"
	srv "github.com/sejamuchhal/taskhub/gateway/server"
)

func main() {
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	fmt.Print("Starting http server")
	server := srv.NewServer(config)

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
