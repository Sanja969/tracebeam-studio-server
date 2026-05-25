FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o tracebeam-server ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/tracebeam-server .

EXPOSE 8080

CMD ["./tracebeam-server"]