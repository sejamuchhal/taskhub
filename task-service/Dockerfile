FROM golang:1.22

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

RUN go build -o /usr/local/bin/app cmd/main.go


CMD ["app", "-b", "0.0.0.0"]
