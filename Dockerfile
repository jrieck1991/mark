FROM golang:latest

WORKDIR /

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o client ./cmd/client
RUN GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

RUN chmod +x client server
