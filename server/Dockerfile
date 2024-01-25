FROM golang:1.21-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.19 AS final

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY start.sh .
COPY app.env .

ENV GIN_MODE=release

EXPOSE 8080

CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
