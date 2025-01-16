FROM golang:1.23.3-alpine

RUN ls -al
RUN apk add --no-cache git gcc musl-dev make sqlite
RUN git clone https://github.com/golang-migrate/migrate.git

WORKDIR ./migrate
RUN ls -al && echo go.mod
RUN cat go.mod

RUN go build -tags 'sqlite3' -ldflags="-X main.Version=$(git describe --tags)" -o $GOPATH/bin/migrate
RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN ls -al

RUN sqlite3 internal/database/shop.db
RUN make migrate_up
#
RUN make build
#
EXPOSE 8080

# build: docker build -t shop_go .
# run:   docker run --rm -p 8080:8080 shop_go

CMD ["./main"]
