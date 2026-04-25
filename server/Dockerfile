FROM golang:1.26-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D -g '' appuser
USER appuser

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]