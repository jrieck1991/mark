FROM golang:latest

RUN mkdir -pv /go/src/github.com/jrieck1991/mark

WORKDIR /go/src/github.com/jrieck1991/mark

COPY . .

RUN GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 go build -o client ./cmd/client/main.go
RUN GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go

RUN chmod +x client server
