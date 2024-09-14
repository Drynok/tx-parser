.PHONY: build up down test fmt lint cover

MAKEFLAGS += --silent

build:
	docker-compose up -d --build
	
up:
	@docker-compose up -d

down:
	docker-compose down

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out