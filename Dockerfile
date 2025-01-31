FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc libc-dev

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main/ .

EXPOSE 8080
CMD ["/app/main"]
