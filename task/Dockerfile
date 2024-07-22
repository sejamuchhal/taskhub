# Stage 1: Build the Go application
FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /app/main cmd/main.go

RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

# Stage 2: Final image
FROM alpine:latest

COPY --from=builder /usr/share/zoneinfo/Asia/Kolkata /usr/share/zoneinfo/Asia/Kolkata
ENV TZ=Asia/Kolkata

WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /go/bin/grpc-health-probe /usr/local/bin/grpc_health_probe

CMD ["./main"]
