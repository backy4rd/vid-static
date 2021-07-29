FROM golang:alpine

WORKDIR /app

RUN apk add ffmpeg

COPY . .

RUN go mod vendor

RUN go build -o app .

CMD ["/app/app"]
