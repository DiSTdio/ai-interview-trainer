FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o app ./cmd

EXPOSE 8080

CMD ["./app"]