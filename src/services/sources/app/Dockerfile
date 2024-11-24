FROM golang:latest

COPY . /runable

WORKDIR /runable/src/services/sources/app

RUN go build main.go

ENTRYPOINT ./main
