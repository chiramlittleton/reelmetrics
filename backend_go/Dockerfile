# Use Golang base image
FROM golang:1.24

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod tidy

# Build the Go app
RUN go build -o main .

# Expose Go backend port
EXPOSE 8002

# Run the Go application
CMD ["/app/main"]
