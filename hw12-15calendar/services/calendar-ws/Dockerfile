## Build
FROM golang:1.20.2-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR  /app

# Copy the code into the container.
COPY services/calendar-ws services/calendar-ws
COPY internal internal

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd services/calendar-ws/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY configs/calendar-ws.yml /configs/
COPY --from=build ["/app/bin/calendar-ws", "/"]

ENTRYPOINT ["/calendar-ws", "-config", "/configs/calendar-ws.yml"]
