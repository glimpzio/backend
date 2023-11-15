FILE="docker/docker-compose.yaml"

.PHONY: up seed down restart run generate

up:
	docker-compose -f $(FILE) up -d

seed:
	docker-compose -f $(FILE) exec -T db psql -U postgres -f /var/data/schema.sql

down:
	docker-compose -f $(FILE) down

restart:
	docker-compose -f $(FILE) restart

run:
	cd src && go run main.go

generate:
	cd src && go get github.com/99designs/gqlgen@v0.17.40 && go run github.com/99designs/gqlgen generate