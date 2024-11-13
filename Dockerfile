FROM golang:1.23.3-alpine

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
