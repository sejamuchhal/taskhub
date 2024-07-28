# Stage 1: Build the Go binary
FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/main cmd/main.go

# Stage 2: Create the final image
FROM alpine:latest

RUN apk update --no-cache && apk add --no-cache tzdata ca-certificates

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]
