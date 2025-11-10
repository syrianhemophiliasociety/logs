.PHONY: build build-dev prod beta

BINARY_NAME=shs-web

build: init
	go build -ldflags="-w -s" -o ${BINARY_NAME} ./cmd/http/main.go

build-dev:
	go build -ldflags="-w -s" -o ${BINARY_NAME} ./cmd/http/main.go

generate: tailwindcss-build
	templ generate -path ./views/

init: htmx-init tailwindcss-init go-init

go-init:
	go mod tidy && \
	go install github.com/a-h/templ/cmd/templ@v0.3.906

htmx-init:
	mkdir -p static/assets/js/htmx && \
	wget https://unpkg.com/hyperscript.org@0.9.14/dist/_hyperscript.min.js -O static/assets/js/htmx/hyperscript.min.js &&\
	wget https://unpkg.com/htmx-ext-json-enc@2.0.2/dist/json-enc.min.js -O static/assets/js/htmx/json-enc.js &&\
	wget https://unpkg.com/htmx-ext-loading-states@2.0.1/dist/loading-states.min.js -O static/assets/js/htmx/loading-states.js &&\
	wget https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js -O static/assets/js/htmx/htmx.min.js

tailwindcss-init:
	mkdir -p static/assets/css &&\
	npm i &&\
	npx @tailwindcss/cli -i static/assets/css/style.css -o static/assets/css/tailwind.css -m

tailwindcss-build:
	npx @tailwindcss/cli -i static/assets/css/style.css -o static/assets/css/tailwind.css

dev:
	air -v > /dev/null
	@if [ $$? != 0 ]; then \
		echo "air was not found, installing it..."; \
		go install github.com/cosmtrek/air@v1.51.0; \
	fi
	air

beta:
	GO_ENV="beta" ./${BINARY_NAME}

prod:
	GO_ENV="prod" ./${BINARY_NAME}

clean:
	go clean

