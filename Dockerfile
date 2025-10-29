# Use the official Go image as a base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application code into the container
COPY . .

# Build the Go application inside the container
RUN make build

# Expose port 443 to the outside world
EXPOSE 443

# Command to run the executable
ENTRYPOINT ["./auth-service"]