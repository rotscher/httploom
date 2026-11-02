.PHONY: generate build test lint tidy run-server run-client

# Regenerate the server and client code from api/openapi.yaml.
generate:
	cd loomserver && go generate ./...
	cd loomclient && go generate ./...

build:
	cd loomserver && go build ./...
	cd loomclient && go build ./...

test:
	cd loomserver && go test ./...
	cd loomclient && go test ./...

lint:
	cd loomserver && golangci-lint run ./...
	cd loomclient && golangci-lint run ./...

tidy:
	cd loomserver && go mod tidy
	cd loomclient && go mod tidy

run-server:
	cd loomserver && go run . -addr :8080

run-client:
	cd loomclient && go run . health
