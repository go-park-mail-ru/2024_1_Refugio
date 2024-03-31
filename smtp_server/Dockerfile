FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main .

EXPOSE 2525

CMD ["./main"]