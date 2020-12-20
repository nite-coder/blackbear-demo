FROM golang:1.15.6 AS builder
WORKDIR /starter
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

WORKDIR /starter/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o starter

FROM alpine:3.12.0
RUN apk update && \
    apk upgrade && \
    apk add --no-cache curl tzdata ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /starter
COPY --from=builder /starter/deployments/database /starter/deployments/database
COPY --from=builder /starter/cmd/starter /starter/starter


RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /starter
USER appuser