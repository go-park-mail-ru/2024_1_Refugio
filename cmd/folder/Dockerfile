FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/folder

EXPOSE 8005

CMD ["./main"]