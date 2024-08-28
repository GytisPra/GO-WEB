# Use a valid Go base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies. Dependencies are cached if the go.mod and go.sum files have not changed
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go application
RUN go build -o myapp ./cmd/go-web

# Expose the port the app runs on
EXPOSE 3000

# Command to run the application
CMD ["./myapp"]