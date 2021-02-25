FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod vendor

RUN go build -o app .

CMD ["/app/app"]
