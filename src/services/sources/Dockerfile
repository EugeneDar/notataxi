FROM golang:1.23 AS builder

WORKDIR /go

COPY . .

RUN go build -o /go/bin/sources src/services/sources/app/main.go

ENTRYPOINT [ "/go/bin/sources" ]
