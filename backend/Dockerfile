FROM golang:1.26.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /doc-hub-api ./cmd/api

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /doc-hub-api /usr/local/bin/doc-hub-api

EXPOSE 8080

CMD ["doc-hub-api"]
