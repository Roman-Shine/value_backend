.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans app