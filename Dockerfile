# Step 1: Build the binary
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

ENV GOARCH=amd64
ENV GOOS=linux
# CGO has to be disabled for alpine
ENV CGO_ENABLED=0 

# Build the Go app
RUN go build -o main ./cmd/main.go

# Step 2: Use a minimal base image to run the application
FROM alpine:latest  

RUN apk --no-cache add ca-certificates bash

# WORKDIR /root/
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main ./

# Ensure the binary is executable
RUN chmod +x main

# Expose port 8080 to the outside world
EXPOSE 8080

ENV DB_HOST=
ENV DB_PORT=
ENV DB_USER=
ENV DB_PASSWORD=
ENV DB_NAME=
ENV DB_PARAMS=sslmode=disable
ENV JWT_SECRET=
ENV BCRYPT_SALT=

# Command to run the executable
CMD ["/app/main"]
