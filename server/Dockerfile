FROM golang:1.26.1 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sentinel .

FROM alpine:latest

WORKDIR /root/