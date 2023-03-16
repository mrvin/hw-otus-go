## Build
FROM golang:1.20.2-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR  /app

# Copy the code into the container.
COPY services/scheduler services/scheduler
COPY internal internal

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd services/scheduler/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY configs/scheduler.yml /configs/
COPY --from=build ["/app/bin/scheduler", "/"]

ENTRYPOINT ["/scheduler", "-config", "/configs/scheduler.yml"]
