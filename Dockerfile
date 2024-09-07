# Stage 1: Build the application
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o acme2kong cmd/acme2kong/main.go

# Stage 2: Create a minimal image to run the application
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/acme2kong .

CMD ["./acme2kong"]