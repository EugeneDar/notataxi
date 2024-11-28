FROM golang:1.23 AS builder

WORKDIR /go

COPY . .

RUN apt-get update && \
    apt-get install wget postgresql-client --yes && \
    mkdir --parents ~/.postgresql && \
    wget "https://storage.yandexcloud.net/cloud-certs/CA.pem" \
         --output-document ~/.postgresql/root.crt && \
    chmod 0655 ~/.postgresql/root.crt

RUN go build -o /go/bin/orders cmd/orders/main.go

ENTRYPOINT [ "/go/bin/orders" ]
