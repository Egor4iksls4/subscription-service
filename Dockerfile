FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest

RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/config ./config
COPY --from=builder /app/.env .

COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appuser .

USER appuser

EXPOSE 8080

CMD ["./main"]