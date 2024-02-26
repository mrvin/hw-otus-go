FROM golang:1.22.0-alpine

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR /calendar/app

COPY cmd/calendar cmd/calendar
COPY Makefile .
COPY internal internal
COPY pkg pkg

COPY deployments/start.sh /

COPY go.mod .
COPY go.sum .
# Download necessary Go modules
RUN go mod download

RUN cd cmd/calendar/ && make build

EXPOSE 8088 55555

CMD ["/start.sh"]

