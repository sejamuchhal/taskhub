.PHONY: all docker gen test test-auth test-gateway test-notification test-task

# Default command on `make`
all: docker

# Build and run all containers
build:
	docker compose up --build

# Clean up docker containers
clean:
	docker compose down

# Generate protobuf files for auth, task and event
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


# Run tests for all services
test: test-auth test-gateway test-notification test-task

# Run tests for the auth service
test-auth:
	go test `go list ./auth/... | grep -v ./auth/pb` -coverprofile=auth/coverage.out

	$(eval coverage=$(shell go tool cover -func=auth/coverage.out | grep total | grep -o -E '[0-9]+' | head -n 1))
	@echo "Total coverage for Auth package $(coverage)%"
	@if [ $(coverage) -lt 60 ]; then \
			echo "Need at least 60% coverage for Auth package"; \
			exit 1; \
	fi

# Run tests for the gateway service
test-gateway:
	go test `go list ./gateway/... | grep -v ./gateway/pb` -coverprofile=gateway/coverage.out
	$(eval coverage=$(shell go tool cover -func=gateway/coverage.out | grep total | grep -o -E '[0-9]+' | head -n 1))
	@echo "Total coverage for Gateway package $(coverage)%"
	@if [ $(coverage) -lt 60 ]; then \
			echo "Need at least 60% coverage for Gateway package"; \
			exit 1; \
	fi

# Run tests for the notification service
test-notification:
	go test `go list ./notification/... | grep -v ./notification/pb` -coverprofile=notification/coverage.out
	$(eval coverage=$(shell go tool cover -func=notification/coverage.out | grep total | grep -o -E '[0-9]+' | head -n 1))
	@echo "Total coverage for Notification package $(coverage)%"
	@if [ $(coverage) -lt 60 ]; then \
			echo "Need at least 60% coverage for Notification package"; \
			exit 1; \
	fi

# Run tests for the task service
test-task:
	go test `go list ./task/... | grep -v ./task/pb` -coverprofile=task/coverage.out
	$(eval coverage=$(shell go tool cover -func=task/coverage.out | grep total | grep -o -E '[0-9]+' | head -n 1))
	@echo "Total coverage for Task package $(coverage)%"
	@if [ $(coverage) -lt 60 ]; then \
			echo "Need at least 60% coverage for Task package"; \
			exit 1; \
	fi

# Run coverage for all services
cover: cover-auth cover-gateway cover-notification cover-task

# Run coverage for the auth service
cover-auth:
	go tool cover -func=auth/coverage.out > auth/coverage.txt

# Run coverage for the gateway service
cover-gateway:
	go tool cover -func=gateway/coverage.out > gateway/coverage.txt

# Run coverage for the notification service
cover-notification:
	go tool cover -func=notification/coverage.out > notification/coverage.txt

# Run coverage for the task service
cover-task:
	go tool cover -func=task/coverage.out > task/coverage.txt
