# Use official Go image as a build environment
FROM golang:1.23 AS build

WORKDIR /app

# Copy backend code
COPY backend . 

# Download dependencies
RUN go mod download
RUN go mod tidy

# Build the Go app
RUN go build -o main cmd/main.go

# Use Ubuntu as runtime (has glibc support)
FROM ubuntu:22.04
WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the binary from build stage
COPY --from=build /app/main .

# Copy frontend files  
COPY frontend ./frontend

# Expose the port your app listens on
EXPOSE 8080

# Run the app
CMD ["./main"]
