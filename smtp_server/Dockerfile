FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN sudo go build -o main .

EXPOSE 587

CMD ["sudo", "./main"]