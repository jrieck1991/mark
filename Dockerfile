FROM golang:latest

WORKDIR /

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o app .

RUN chmod +x app

ENTRYPOINT [ "./app" ]