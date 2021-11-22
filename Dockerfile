FROM golang:1.17.3

ENV CC=arm-linux-gnueabi-gcc \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=arm \
    GOARM=5

RUN apt-get update && apt-get install -y gcc-arm-linux-gnueabi

WORKDIR /app

COPY go.* .
RUN go mod download
