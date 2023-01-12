FROM golang:1.19-alpine3.16 AS builder

RUN mkdir /app
WORKDIR /app

RUN apk add git bash make

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build

CMD ["/app/bookline"]

FROM alpine:3.14
RUN apk add --update --no-cache ca-certificates tzdata bash curl
SHELL ["/bin/bash", "-c"]

# set up nsswitch.conf for Go's "netgo" implementation
# https://github.com/gliderlabs/docker-alpine/issues/367#issuecomment-424546457
RUN test ! -e /etc/nsswitch.conf && echo 'hosts: files dns' > /etc/nsswitch.conf

COPY --from=builder /app/bookline /usr/local/bin/

CMD ["/usr/local/bin/bookline"]