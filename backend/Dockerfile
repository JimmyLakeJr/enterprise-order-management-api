FROM golang:1.25.5-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:3.23

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /app/server /app/server

USER app

EXPOSE 8080

CMD ["/app/server"]
