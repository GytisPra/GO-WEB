# Use the official Go image as a base image
FROM golang:1.21-alpine

# Install debugging tools
RUN apk add --no-cache tree

# Set the working directory inside the container
WORKDIR /app

# Debug: Print current directory
RUN pwd

# Copy the Go modules and dependencies files
COPY go.mod go.sum ./

# Debug: List contents after copying go.mod and go.sum
RUN echo "Contents after copying go.mod and go.sum:" && ls -la

# Download dependencies. Dependencies are cached if the go.mod and go.sum files have not changed
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

# Copy the entire project
COPY . .

# Debug: List contents of /app
RUN echo "Contents of /app:" && ls -la /app

# Debug: Show directory structure
RUN echo "Directory structure:" && tree /app

# Debug: Print Go version
RUN go version

# Debug: Print go env
RUN go env

# Build the Go application
# Adjust the path to where your main.go file is located
RUN go build -v -o myapp ./cmd/go-web

# Expose the port on which the app will run
EXPOSE 3000

# Command to run the application
CMD ["./myapp"]