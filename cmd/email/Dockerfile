FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/email

EXPOSE 8002

CMD ["./main"]