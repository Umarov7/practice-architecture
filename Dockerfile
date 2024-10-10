FROM golang:1.23.2 AS builder

WORKDIR /app

COPY . .
COPY .env .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../myapp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./myapp"]