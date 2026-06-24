FROM golang:1.25.5-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

FROM alpine:3.23

WORKDIR /app
COPY --from=builder /app/server /app/server
RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["/app/server"]
