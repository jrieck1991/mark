.PHONY: run clean

APP_NAME=mark

run: clean build
	docker-compose up -d

clean:
	docker-compose down

build:
	docker-compose build
