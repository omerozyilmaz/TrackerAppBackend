FROM golang:1.19-alpine AS builder

WORKDIR /app

# Go proxy ayarı ekleyin
ENV GOPROXY=https://proxy.golang.org,direct

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download -x  # -x bayrağı daha fazla hata ayıklama bilgisi sağlar

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller image for the final application
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your application runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"] 