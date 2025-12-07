FROM golang:1.25-alpine AS builder

RUN apk add --no-cache chromium

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /arbeit-scraper .

FROM alpine:latest
WORKDIR /root/

RUN apk add --no-cache chromium

COPY --from=builder /arbeit-scraper .

EXPOSE 8080
CMD ["./arbeit-scraper"]