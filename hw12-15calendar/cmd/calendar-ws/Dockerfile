## Build
FROM golang:1.20.2-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make
RUN apk --no-cache add tzdata

WORKDIR  /app

# Copy the code into the container.
COPY cmd/calendar-ws cmd/calendar-ws
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/calendar-ws/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/app/bin/calendar-ws", "/"]
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE 8080

ENTRYPOINT ["/calendar-ws"]
