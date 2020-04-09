.PHONY: run clean

APP_NAME=mark

run: build
	docker-compose up -d

clean:
	docker-compose down

build:
	docker build -t $(APP_NAME) .
