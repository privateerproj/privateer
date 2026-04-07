FROM golang:1.26.1-alpine3.22@sha256:07e91d24f6330432729082bb580983181809e0a48f0f38ecde26868d4568c6ac AS build

WORKDIR /app

COPY . .

RUN apk add --no-cache make git && \
  go mod tidy && \
  make build

FROM alpine:3.22@sha256:55ae5d250caebc548793f321534bc6a8ef1d116f334f18f4ada1b2daad3251b2

WORKDIR /app

COPY --from=build /app/pvtr /app/pvtr

CMD ["/app/pvtr", "--help"]
