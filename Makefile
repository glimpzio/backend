FILE="docker/docker-compose.yaml"

.PHONY: up down

up:
	docker-compose -f $(FILE) up -d 

down:
	docker-compose -f $(FILE) down

restart:
	docker-compose -f $(FILE) restart