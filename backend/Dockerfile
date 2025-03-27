FROM golang:1.24-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Verify the binary exists and is executable
RUN ls -la main && chmod +x main

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/main"] 