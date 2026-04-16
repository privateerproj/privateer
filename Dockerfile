FROM golang:1.26.2-alpine3.22@sha256:c259ff7ffa06f1fd161a6abfa026573cf00f64cfd959c6d2a9d43e3ff63e8729 AS build

WORKDIR /app

COPY . .

RUN apk add --no-cache make git && \
  go mod tidy && \
  make build

FROM alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

WORKDIR /app

COPY --from=build /app/pvtr /app/pvtr

CMD ["/app/pvtr", "--help"]
