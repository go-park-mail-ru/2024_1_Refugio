FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/auth-vk

EXPOSE 8007

CMD ["./main"]