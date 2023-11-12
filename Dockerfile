#Build stage
FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

#Final stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .

# cpy the app.env file to the container
COPY --from=builder /app/app.env .

EXPOSE 8080

CMD ["/app/main"]
