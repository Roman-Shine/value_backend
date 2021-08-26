.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	set GOOS=linux && go mod download && go build -o ./.bin/app ./cmd/app/main.go
	#go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans app