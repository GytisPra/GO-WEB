# Use the official Go image as a base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies files
COPY go.mod go.sum ./

# Download dependencies. Dependencies are cached if the go.mod and go.sum files have not changed
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

# Copy the source code to the container
COPY . .

# Build the Go application
RUN cd ./cmd/go-web && go build -o myapp

# Expose the port on which the app will run
EXPOSE 3000

# Command to run the application
CMD ["./cmd/go-web/myapp"]