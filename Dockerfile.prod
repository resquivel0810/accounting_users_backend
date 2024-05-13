# Use an official Golang runtime as a parent image
FROM golang:alpine

# Install git, required for fetching Go dependencies
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project and preserve the structure
COPY . .

# Download dependencies using go mod
RUN go mod download

# Build the application to produce a binary
RUN go build -o app ./cmd/api/*.go

# Command to run the executable
CMD ["./app"]
