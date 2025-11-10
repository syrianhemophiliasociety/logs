FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN apk add --no-cache make npm nodejs

RUN make init
RUN make generate build

FROM alpine:latest AS run

WORKDIR /app
COPY --from=build /app/shs-web ./shs-web

EXPOSE 3003

CMD ["./shs-web"]

