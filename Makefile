.PHONY: all build shs-server

SERVER_BINARY_NAME=shs-server
MIGRATOR_BINARY_NAME=shs-migrator

all: build-server build-migrator

build: init build-server build-migrator

build-server: generate
	go build -ldflags="-w -s" -o ${SERVER_BINARY_NAME} ./cmd/http/main.go

build-migrator: build-server
	go build -ldflags="-w -s" -o ${MIGRATOR_BINARY_NAME} ./cmd/migrator/main.go

init: htmx-init tailwindcss-init go-init

migrate: build-migrator
	./${MIGRATOR_BINARY_NAME}

generate:
	templ generate -path ./web/views/

go-init:
	go mod tidy && \
	go install github.com/a-h/templ/cmd/templ@v0.3.906

htmx-init:
	mkdir -p web/static/assets/js/htmx && \
	wget https://unpkg.com/hyperscript.org@0.9.14/dist/_hyperscript.min.js -O web/static/assets/js/htmx/hyperscript.min.js &&\
	wget https://unpkg.com/htmx-ext-json-enc@2.0.2/dist/json-enc.min.js -O web/static/assets/js/htmx/json-enc.js &&\
	wget https://unpkg.com/htmx-ext-loading-states@2.0.1/dist/loading-states.min.js -O web/static/assets/js/htmx/loading-states.js &&\
	wget https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js -O web/static/assets/js/htmx/htmx.min.js

tailwindcss-init:
	mkdir -p web/static/assets/css &&\
	npm i &&\
	npx @tailwindcss/cli -i web/static/assets/css/style.css -o web/static/assets/css/tailwind.css -m

tailwindcss-build:
	npx @tailwindcss/cli -i web/static/assets/css/style.css -o web/static/assets/css/tailwind.css

tailwindcss-server:
	npx @tailwindcss/cli -i web/static/assets/css/style.css -o web/static/assets/css/tailwind.css --watch

dev:
	air -v > /dev/null
	@if [ $$? != 0 ]; then \
		echo "air was not found, installing it..."; \
		go install github.com/cosmtrek/air@v1.51.0; \
	fi
	export `cat .env | xargs` && air

dev-test:
	air -v > /dev/null
	@if [ $$? != 0 ]; then \
		echo "air was not found, installing it..."; \
		go install github.com/cosmtrek/air@v1.51.0; \
	fi
	export `cat .env.ci | xargs` && air

shs-server:
	./${MIGRATOR_BINARY_NAME} &&\
	./${SERVER_BINARY_NAME}

test:
	@npx playwright test --reporter=list

test-ci: build
	@npx playwright test --reporter=dot

clean:
	go clean
