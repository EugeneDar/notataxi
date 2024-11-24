FROM golang:1.23 AS builder

WORKDIR /go

COPY . .

RUN go build -o /go/bin/orders src/services/orders/app/main.go

ENTRYPOINT [ "/go/bin/orders" ]
