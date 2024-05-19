FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/smtp

EXPOSE 587

CMD ["./main"]