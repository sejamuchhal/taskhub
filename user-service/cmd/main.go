package main

import (
	"fmt"

	srv "github.com/sejamuchhal/taskhub/user-service/internal/server"
)

func main() {
	fmt.Print("Starting http server")
	server := srv.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
