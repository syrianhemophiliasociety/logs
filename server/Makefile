.PHONY: all build shs-server

SERVER_BINARY_NAME=shs-server
MIGRATOR_BINARY_NAME=shs-migrator

all: build-server build-migrator

build: init generate build-server build-migrator

build-server: init
	go build -ldflags="-w -s" -o ${SERVER_BINARY_NAME} ./cmd/http/main.go

build-migrator: init
	go build -ldflags="-w -s" -o ${MIGRATOR_BINARY_NAME} ./cmd/migrator/main.go

init:
	go mod tidy

migrate: build-migrator
	./${MIGRATOR_BINARY_NAME}

dev:
	air -v > /dev/null
	@if [ $$? != 0 ]; then \
		echo "air was not found, installing it..."; \
		go install github.com/cosmtrek/air@v1.51.0; \
	fi

	air

shs-server:
	./${MIGRATOR_BINARY_NAME} &&\
	./${SERVER_BINARY_NAME}

clean:
	go clean
