FROM golang:1.23 AS builder

WORKDIR /go

COPY . .

RUN apt-get update && \
    apt-get install wget postgresql-client --yes && \
    mkdir --parents /go/.postgresql && \
    wget "https://storage.yandexcloud.net/cloud-certs/CA.pem" \
         --output-document /go/.postgresql/root.crt && \
    chmod 0655 /go/.postgresql/root.crt

RUN go build -o /go/bin/orders cmd/orders/main.go

ENTRYPOINT [ "/go/bin/orders" ]
