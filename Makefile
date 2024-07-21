PROTOC := protoc
PROTO_SRC := protos/task.proto
PROTO_PATH := protos
TASK_SERVICE_OUT := task-service/pb
API_GATEWAY_OUT := api-gateway/client


all: docker

docker:
	docker compose up --build

clean:
	rm -rf task-service/pb/*.go
	rm -rf api-gateway/client/*.go

gen:
	protoc -I=protos \
		--go_out=task-service/pb/task --go_opt=paths=source_relative \
		--go-grpc_out=task-service/pb/task --go-grpc_opt=paths=source_relative \
		protos/task.proto
	protoc -I=protos \
		--go_out=api-gateway/client --go_opt=paths=source_relative \
		--go-grpc_out=api-gateway/client --go-grpc_opt=paths=source_relative \
		protos/task.proto
	protoc -I=protos \
		--go_out=task-service/pb/event --go_opt=paths=source_relative \
		--go-grpc_out=task-service/pb/event --go-grpc_opt=paths=source_relative \
		protos/event.proto
	protoc -I=protos \
		--go_out=notification-service/pb --go_opt=paths=source_relative \
		--go-grpc_out=notification-service/pb --go-grpc_opt=paths=source_relative \
		protos/event.proto
	protoc -I=protos \
		--go_out=user-service/client/task --go_opt=paths=source_relative \
		--go-grpc_out=user-service/client/task --go-grpc_opt=paths=source_relative \
		protos/task.proto

task-service/pb:
	mkdir -p task-service/pb

api-gateway/client:
	mkdir -p api-gateway/client

docker:
	docker compose up --build

# pg-ip:
#	docker ps
# 	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' taskhub-db-1
# 	psql -U ${DB_USERNAME} -d task_db