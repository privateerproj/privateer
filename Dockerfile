FROM golang:1.26.4-alpine3.22@sha256:727cfc3c40be55cd1bc9a4a059406b28a059857e3be752aa9d09531e12c20c56 AS build

WORKDIR /app

COPY . .

RUN apk add --no-cache make git && \
  go mod tidy && \
  make build

FROM alpine:3.24@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4

WORKDIR /app

COPY --from=build /app/pvtr /app/pvtr

CMD ["/app/pvtr", "--help"]
