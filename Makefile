all: docker

docker:
	docker compose up --build

gen:
	protoc -I=protos \
		--go_out=task/pb/task --go_opt=paths=source_relative \
		--go-grpc_out=task/pb/task --go-grpc_opt=paths=source_relative \
		protos/task.proto
	protoc -I=protos \
		--go_out=gateway/pb/task --go_opt=paths=source_relative \
		--go-grpc_out=gateway/pb/task --go-grpc_opt=paths=source_relative \
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
		--go_out=auth/pb --go_opt=paths=source_relative \
		--go-grpc_out=auth/pb --go-grpc_opt=paths=source_relative \
		protos/auth.proto
	protoc -I=protos \
		--go_out=gateway/pb/auth --go_opt=paths=source_relative \
		--go-grpc_out=gateway/pb/auth --go-grpc_opt=paths=source_relative \
		protos/auth.proto


docker:
	docker compose up --build