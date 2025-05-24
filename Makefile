.DEFAULT_GOAL:=help
.PHONY: help

APP_NAME:=customer-api
DOCKER_IMAGE_NAME:=wishlist-customer-api
PWD=$(shell pwd)
PORT=8080

APP_DIR=/app


start_docker: 
	@docker-compose --progress quiet up -d --no-build --quiet-pull

docker-run: start_docker
	@docker exec -it $(APP_NAME) sh -c "$(CMD)"

generate:
	@echo 'generating files ...'
	@make docker-run CMD="go generate ./..."
	@echo 'done!'

test: generate
	@echo 'generating'
	@make docker-run CMD='go test -coverprofile=coverage.out ./... ; go tool cover -func=coverage.out; go tool cover -html=coverage.out -o coverage.html'