###############
# Build stage #
###############
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o txparser ./cmd/main.go

#############
# Run stage #
#############
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/txparser .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./txparser"]
