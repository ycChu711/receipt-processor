# Build stage
FROM golang:1.19-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o receipt-processor .

# Final stage - minimal runtime image
FROM alpine:3.17

# Add CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/receipt-processor .

# Expose port
EXPOSE 8080

# Command to run
ENTRYPOINT ["./receipt-processor"]