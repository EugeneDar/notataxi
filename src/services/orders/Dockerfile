FROM golang:1.23 AS builder
WORKDIR /app

COPY app .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath -mod=readonly -o /go/bin/app

ENTRYPOINT [ "/go/bin/app" ]
