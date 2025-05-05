FROM golang:1.24

WORKDIR /usr/src/app

COPY go.mod go.sum ./

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/main.go

CMD ["app"]