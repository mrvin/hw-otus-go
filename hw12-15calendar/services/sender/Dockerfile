## Build
FROM golang:1.19.2-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

WORKDIR  /app

# Copy the code into the container.
COPY services/sender services/sender
COPY internal internal

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd services/sender/ && go build -o sender

## Deploy
FROM alpine:3.16.2

WORKDIR /

COPY configs/sender.yml /configs/
COPY --from=build ["/app/services/sender/sender", "/"]

ENTRYPOINT ["/sender", "-config", "/configs/sender.yml"]
