FILE="docker/docker-compose.yaml"

.PHONY: up down restart run generate

up:
	docker-compose -f $(FILE) up -d 

down:
	docker-compose -f $(FILE) down

restart:
	docker-compose -f $(FILE) restart

run:
	go run src/main.go

generate:
	cd src && go run github.com/99designs/gqlgen generate