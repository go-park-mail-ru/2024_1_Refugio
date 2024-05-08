FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main .

EXPOSE 587

CMD ["./main"]