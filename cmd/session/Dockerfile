FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/session

EXPOSE 8003

CMD ["./main"]