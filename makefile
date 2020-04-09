.PHONY: run clean

run: clean
	docker-compose build
	docker-compose up -d

clean:
	docker-compose down