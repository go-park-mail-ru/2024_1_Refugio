FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main .

EXPOSE 8082

CMD ["./main"]