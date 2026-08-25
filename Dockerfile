FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o taskflow-api ./cmd/api

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/taskflow-api .

EXPOSE 8080

CMD ["./taskflow-api"]