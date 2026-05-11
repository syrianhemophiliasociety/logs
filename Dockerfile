FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go install github.com/a-h/templ/cmd/templ@v0.2.731 &&\
    apk add --no-cache make npm nodejs

RUN make init &&\
    make build-server &&\
    make build-migrator

FROM alpine:latest AS run

RUN apk add --no-cache make

WORKDIR /app
COPY --from=build /app/shs-logs-server ./shs-logs-server
COPY --from=build /app/shs-logs-migrator ./shs-logs-migrator
COPY --from=build /app/Makefile ./Makefile

EXPOSE 3000

CMD ["make", "shs-logs-server"]
