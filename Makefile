.PHONY:
.SILENT:
.DEFAULT_GOAL := run

IS_TODO_BACKEND_RUNNING := $(shell docker ps --filter name=todo_backend --filter status=running -aq)
IS_TODO_BACKEND_EXITED := $(shell docker ps --filter name=todo_backend -aq)
IS_TODO_BACKEND_IMAGE := $(shell docker images --filter=reference="todo_backend" -aq)

IS_TODO_DATABASE_RUNNING := $(shell docker ps --filter name=todo_database --filter status=running -aq)
IS_TODO_DATABASE_EXITED := $(shell docker ps --filter name=todo_database -aq)

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: clean
	docker-compose up --build --detach

stop:
	docker-compose down

clean:
ifneq ($(strip $(IS_TODO_BACKEND_RUNNING)),)
	docker stop $(IS_TODO_BACKEND_RUNNING)
endif

ifneq ($(strip $(IS_TODO_BACKEND_EXITED)),)
	docker rm $(IS_TODO_BACKEND_EXITED)
endif

ifneq ($(strip $(IS_TODO_BACKEND_IMAGE)),)
	docker rmi $(IS_TODO_BACKEND_IMAGE)
endif

ifneq ($(strip $(IS_TODO_DATABASE_RUNNING)),)
	docker stop $(IS_TODO_DATABASE_RUNNING)
endif

ifneq ($(strip $(IS_TODO_DATABASE_EXITED)),)
	docker rm $(IS_TODO_DATABASE_EXITED)
endif

cover:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html

swag:
	swag init -g ./cmd/app/main.go