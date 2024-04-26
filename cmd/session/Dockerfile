FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/mail

EXPOSE 8003

CMD ["./main"]