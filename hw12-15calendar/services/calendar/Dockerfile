FROM golang:1.20.2-alpine

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR /calendar/app

COPY services/calendar services/calendar
COPY Makefile .
COPY internal internal
COPY pkg pkg

COPY deployments/start.sh /

COPY go.mod .
COPY go.sum .
# Download necessary Go modules
RUN go mod download

RUN cd services/calendar/ && make build

EXPOSE 8088 55555

CMD ["/start.sh"]

