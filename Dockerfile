FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go build -o bsky-stats-api .

EXPOSE 8080

CMD ["./bsky-stats-api"]