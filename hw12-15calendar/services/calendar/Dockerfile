FROM golang:1.19.2-alpine

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk --update add build-base

WORKDIR /calendar/app

COPY services/calendar services/calendar
COPY Makefile .
COPY internal internal

COPY configs/calendar.yml configs/

COPY deployments/start.sh /

COPY go.mod .
COPY go.sum .
# Download necessary Go modules
RUN go mod download

RUN cd services/calendar/ && make build

EXPOSE 8080

CMD ["/start.sh"]

