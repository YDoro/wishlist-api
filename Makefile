.DEFAULT_GOAL:=help
.PHONY: help

APP_NAME:=customer-api
DOCKER_IMAGE_NAME:=wishlist-customer-api
PWD=$(shell pwd)
PORT=8080

APP_DIR=/app


_force-build: ## Force build the docker image
	@docker compose build --no-cache

_setEnv:
	@if [ -f .env ]; then \
		sed -i '/^ENV=/d' .env; \
	fi
	@echo "ENV=$(VAL)" >> .env

start-dev: _force-build
	@make _setEnv VAL=dev
	@docker compose up -d

start: _force-build
	@make _setEnv VAL=prod
	@docker compose up

stop:
	@docker-compose --progress quiet down

start_docker: 
	@docker-compose --progress quiet up -d --no-build --quiet-pull

docker-run:
	@docker exec -it $(APP_NAME) sh -c "$(CMD)"

generate:
	@echo 'generating files ...'
	@make docker-run CMD="go generate ./..."
	@echo 'done!'

test: generate
	@echo 'generating'
	@make docker-run CMD='go test -coverprofile=coverage.out ./internal/... ; go tool cover -func=coverage.out; go tool cover -html=coverage.out -o coverage.html'

docsx:
	@echo 'generating docs ...'
	@make docker-run CMD='swag init -g ./cmd/main.go -o ./docs'
	@echo 'done!'