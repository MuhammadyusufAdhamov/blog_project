FROM golang:1.19.1-alpine3.16 AS builder

WORKDIR /app

COPY . .

RUN apk add curl
RUN go build -o main cmd/main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
 
FROM alpine:3.16

WORKDIR /app
RUN mkdir media

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY migrations ./migrations
COPY templates ./templates

EXPOSE 8007

CMD ["/app/main"]

# docker run --env-file ./.env --name blogApp --network blog-network -p 8000:8000 -d blog:latest
