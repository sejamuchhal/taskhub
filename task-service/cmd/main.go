package taskservice

import (
	_ "github.com/joho/godotenv/autoload"
	pb "github.com/sejamuchhal/task-management/proto"
)

type server struct {
	pb.UnimplementedTaskServiceServer
}

func main() {

}
