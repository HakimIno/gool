# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (needed for go mod download with private repos)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gool main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/gool .

# Make it executable
RUN chmod +x ./gool

# Add to PATH
ENV PATH="/root:${PATH}"

ENTRYPOINT ["./gool"] 