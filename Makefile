.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	./.bin/app

cover:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html


