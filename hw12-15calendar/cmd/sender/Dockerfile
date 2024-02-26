## Build
FROM golang:1.22.0-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR  /app

# Copy the code into the container.
COPY cmd/sender cmd/sender
COPY internal internal

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/sender/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/app/bin/sender", "/"]

ENTRYPOINT ["/sender"]
