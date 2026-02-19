# Official Go base image
FROM golang:1.26-alpine

# Working directory inside the container
WORKDIR /app

# Copy go.mod first (optimizes cache)
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the binary
RUN go build -o server ./cmd/server

# Expose port 8080
EXPOSE 8080

# Command executed when the container starts
CMD ["./server"]
