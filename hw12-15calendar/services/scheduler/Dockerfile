## Build
FROM golang:1.19.2-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

WORKDIR  /app

# Copy the code into the container.
COPY services/scheduler services/scheduler
COPY internal internal

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd services/scheduler/ && go build -o scheduler

## Deploy
FROM alpine:3.16.2

WORKDIR /

COPY configs/scheduler.yml /configs/
COPY --from=build ["/app/services/scheduler/scheduler", "/"]

ENTRYPOINT ["/scheduler", "-config", "/configs/scheduler.yml"]
