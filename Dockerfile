FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build -o app main.go

EXPOSE 80

CMD ["./app"]
