FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o sources ./app/main.go

EXPOSE 9090

CMD ["./sources"]