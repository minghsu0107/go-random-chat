FROM golang:1.17 AS builder

RUN mkdir -p /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/main.go

FROM alpine:3.14

RUN mkdir -p /app/template
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY --from=builder /app/template /app/template
RUN apk add --no-cache bash

ENTRYPOINT ["/app/server"]