FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd/questionnaire

EXPOSE 8006

CMD ["./main"]