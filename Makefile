.PHONY:
.SILENT:
.DEFAULT_GOAL := run

IS_LETS_GO_AGGREGATOR_RUNNING := $(shell docker ps --filter name=todo_backend --filter status=running -aq)
IS_LETS_GO_AGGREGATOR_EXITED := $(shell docker ps --filter name=todo_backend -aq)
IS_LETS_GO_AGGREGATOR := $(shell docker images --filter=reference="*/todo_backend" -aq)

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: clear
	docker-compose up --build --detach

stop:
	docker-compose down

clear: stop

ifneq ($(strip $(IS_TODO_BACKEND_RUNNING)),)
	docker stop $(IS_TODO_BACKEND_RUNNING)
endif

ifneq ($(strip $(IS_TODO_BACKEND_EXITED)),)
	docker rm $(IS_TODO_BACKEND_EXITED)
endif

ifneq ($(strip $(IS_TODO_BACKEND)),)
	docker rmi $(IS_TODO_BACKEND)
endif

cover:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html

swag:
	swag init -g internal/app/app.go