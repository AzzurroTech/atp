# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
# We specify the package as "." because main.go is in the root
RUN CGO_ENABLED=0 GOOS=linux go build -o /atp .

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /atp .
EXPOSE 8080
CMD ["./atp"]