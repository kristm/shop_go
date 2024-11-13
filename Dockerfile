FROM golang:1.23.3-alpine

RUN apk add --no-cache git gcc musl-dev make
RUN git clone https://github.com/golang-migrate/migrate.git
RUN echo $GOPATH

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
#
RUN go build cmd/main.go
#
EXPOSE 8080
#
CMD ["./main"]
