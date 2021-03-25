FROM golang:alpine

WORKDIR /app

COPY . .

RUN wget -O ./static/photos/default-avatar.png https://i.imgur.com/E99FjUR.png

RUN wget -O ./static/photos/default-icon.png https://i.imgur.com/QnEb3Pu.png

RUN go mod vendor

RUN go build -o app .

CMD ["/app/app"]
