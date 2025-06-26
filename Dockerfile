FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o proxy main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/proxy .
EXPOSE 80
CMD ["./proxy"]
