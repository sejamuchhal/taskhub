all: docker

docker:
	docker compose up --build

gen:
	protoc -I=protos \
		--go_out=task/pb/task --go_opt=paths=source_relative \
		--go-grpc_out=task/pb/task --go-grpc_opt=paths=source_relative \
		protos/task.proto
	protoc -I=protos \
		--go_out=api-gateway/client --go_opt=paths=source_relative \
		--go-grpc_out=api-gateway/client --go-grpc_opt=paths=source_relative \
		protos/task.proto
	protoc -I=protos \
		--go_out=task/pb/event --go_opt=paths=source_relative \
		--go-grpc_out=task/pb/event --go-grpc_opt=paths=source_relative \
		protos/event.proto
	protoc -I=protos \
		--go_out=notification/pb --go_opt=paths=source_relative \
		--go-grpc_out=notification/pb --go-grpc_opt=paths=source_relative \
		protos/event.proto
	protoc -I=protos \
		--go_out=auth/client/task --go_opt=paths=source_relative \
		--go-grpc_out=auth/client/task --go-grpc_opt=paths=source_relative \
		protos/task.proto


docker:
	docker compose up --build

# pg-ip:
#	docker ps
# 	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' taskhub-db-1
# 	psql -U ${DB_USERNAME} -d task_db