FROM golang:1.23.1-alpine

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

RUN apk update && apk add make && apk add tzdata

WORKDIR /calendar/app

COPY cmd/calendar cmd/calendar
COPY Makefile .
COPY internal internal
COPY pkg pkg

COPY deployments/start.sh /

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/calendar/ && make build

ENV TZ=Europe/Moscow

EXPOSE 8088 55555

CMD ["/start.sh"]

