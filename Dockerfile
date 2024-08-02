# Start a new stage with a smaller base image
FROM alpine:latest

# Set the working directory to /app in the second stage
WORKDIR /app

# Install curl. --no-cache prevents pkg repo index caching.
# This helps reducing the docker image size
# AWS ECS requires curl to monitor health of golang app

RUN apk add --no-cache curl

# Copy the binary from the build stage to the second stage
COPY server-linux /app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server-linux"]

