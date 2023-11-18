FILE="docker/docker-compose.yaml"

.PHONY: up seed down run generate deploy

up:
	docker-compose -f $(FILE) up -d

seed:
	docker-compose -f $(FILE) exec -T db psql -U postgres -f /var/data/schema.sql

down:
	docker-compose -f $(FILE) down

run:
	cd src && go run main.go

generate:
	cd src && go get github.com/99designs/gqlgen@v0.17.40 && go run github.com/99designs/gqlgen generate

deploy:
	cd aws && cdk deploy