FROM golang:1.26.1-alpine3.22@sha256:07e91d24f6330432729082bb580983181809e0a48f0f38ecde26868d4568c6ac AS build

WORKDIR /app

COPY . .

RUN apk add --no-cache make git && \
  go mod tidy && \
  make build

FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

WORKDIR /app

COPY --from=build /app/pvtr /app/pvtr

CMD ["/app/pvtr", "--help"]
