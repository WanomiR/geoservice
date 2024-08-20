FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod init geoservice && go mod tidy

RUN go build -o main ./cmd/api


FROM alpine:latest

COPY --from=builder /app/main /main
COPY --from=builder /app/.env /.env
COPY --from=builder /app/docs /docs

EXPOSE 8888

CMD ["/main"]