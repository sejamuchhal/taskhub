package main

import (
	"fmt"

	"github.com/sejamuchhal/task-management/user-service/internal/server"
)

func main() {
	fmt.Print("Starting http server")
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
