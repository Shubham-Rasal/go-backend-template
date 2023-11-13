#Build stage
FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build \
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main main.go
# RUN go build -o main main.go

RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin/

#Final stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/migrate ./migrate
COPY --from=builder /app/start.sh .
RUN chmod +x start.sh

# cpy the app.env file to the container
# COPY --from=builder /app/app.env .
COPY --from=builder /app/db/migration ./migration

EXPOSE 8080

CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]
