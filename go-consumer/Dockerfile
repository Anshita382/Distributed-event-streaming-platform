FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer ./cmd/consumer

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/consumer /app/consumer
EXPOSE 2112
ENTRYPOINT ["/app/consumer"]
