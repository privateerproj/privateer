FROM golang:1.23.4-alpine3.21 AS build

WORKDIR /app

COPY . .

RUN apk add --no-cache make git && \
  go mod tidy && \
  make build

FROM alpine:3.21

WORKDIR /app

COPY --from=build /app/privateer /app/privateer

CMD ["/app/privateer", "--help"]
