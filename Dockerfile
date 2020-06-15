#
# Builder
#
FROM golang:alpine3.12 as builder

COPY go.mod go.sum /opt/

WORKDIR /opt

RUN go mod download

ADD . /opt

RUN go build -o app

#
# Runtime
#
FROM alpine:3.12

WORKDIR /opt

COPY --from=builder /opt/app .

ENTRYPOINT ["/opt/app"]
