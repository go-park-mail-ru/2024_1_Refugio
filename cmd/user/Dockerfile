FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/user

EXPOSE 8001

CMD ["./main"]