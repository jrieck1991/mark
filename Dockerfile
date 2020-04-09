FROM golang:latest

WORKDIR /

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o client ./cmd/client/main.go
RUN GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go

RUN chmod +x client server
