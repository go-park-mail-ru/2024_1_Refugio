FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/auth

EXPOSE 8004

CMD ["./main"]