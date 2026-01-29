# Use an official Golang runtime as a parent image
FROM golang:alpine

# Install git and make, required for fetching Go dependencies and running make commands
RUN apk add --no-cache git make

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project and preserve the structure
COPY . .

# Download dependencies using go mod
RUN go mod download

# Install swag for generating Swagger documentation (usar misma versión que en go.mod)
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.1

# Add Go bin to PATH for swag command (GOPATH default is /go in golang:alpine)
ENV PATH=$PATH:/go/bin

# Generate Swagger documentation (creates backend/docs package)
RUN swag init -g cmd/api/main.go -o docs

# Build the application to produce a binary
RUN go build -o app ./cmd/api/*.go

# Command to run the executable
CMD ["./app"]
