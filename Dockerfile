FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go build -o server main.go

FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]