.DEFAULT_GOAL:=help
.PHONY: help

APP_NAME:=customer-api
DOCKER_IMAGE_NAME:=wishlist-customer-api
PWD=$(shell pwd)
PORT=8080

APP_DIR=/app
#DB_STRING=postgres://$$$$DB_USER:$$$$DB_PASSWORD@$$$$DB_HOST:$$$$DB_PORT/$$$$DB_NAME?sslmode=disable

	
docker-run:
	@docker exec -it $(APP_NAME) sh -c "$(DOCKER_COMMAND)"

test:
	@make docker-run DOCKER_COMMAND='go test ./...'
#migrate:
#	@make docker-run DOCKER_COMMAND='migrate -path ./internal/customer/infra/db/postgres/migrations -database "$(DB_STRING)" up'
