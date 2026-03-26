# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /atp ./cmd/server

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /atp .
EXPOSE 8080
CMD ["./atp"]