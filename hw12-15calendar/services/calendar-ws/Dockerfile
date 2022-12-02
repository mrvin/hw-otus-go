## Build
FROM golang:1.19.2-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

WORKDIR  /app

# Copy the code into the container.
COPY services/calendar-ws services/calendar-ws
COPY internal internal

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd services/calendar-ws/ && go build -o calendar-ws

## Deploy
FROM alpine:3.16.2

WORKDIR /

COPY configs/calendar-ws.yml /configs/
COPY --from=build ["/app/services/calendar-ws/calendar-ws", "/"]

ENTRYPOINT ["/calendar-ws", "-config", "/configs/calendar-ws.yml"]
