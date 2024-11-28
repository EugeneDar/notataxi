FROM golang:1.23 AS builder

WORKDIR /go

COPY . .

RUN go build -o /go/bin/sources cmd/sources/main.go

ENTRYPOINT [ "/go/bin/sources" ]
