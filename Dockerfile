FROM golang:1.19.3-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.16
WORKDIR /app

COPY --from=builder /app/main .
COPY db/migration ./db/migration
COPY app.env .

EXPOSE 8000
CMD ["/app/main"]