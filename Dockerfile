FROM golang:1.19-alpine3.16

RUN mkdir /app
WORKDIR /app

RUN apk add git bash

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v ./...

CMD ["/app/bookline"]