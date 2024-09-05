FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod init auth && go mod tidy

RUN go build -o main ./cmd/api


FROM alpine:latest

COPY --from=builder /app/main /main

EXPOSE 2000

CMD ["/main"]
