# Start from the official Golang image to build our application.
FROM golang:1.20 as builder

WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dice .
FROM alpine:latest

# Install ca-certificates in case your application makes outgoing HTTPS requests.
RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /app/dice .

# Expose port 8080 to the outside world.
EXPOSE 8080

# Command to run the executable.
CMD ["./dice"]
