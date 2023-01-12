FROM golang:1.19-alpine3.16

RUN mkdir /app
WORKDIR /app

RUN apk add git bash make

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build

CMD ["/app/bookline"]