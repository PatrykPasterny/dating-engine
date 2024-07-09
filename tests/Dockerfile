# Use an official Golang runtime as a parent image
FROM golang:1.22-alpine

# Set the working directory to /tmp/build/tests
WORKDIR /tmp/build/tests

# Create a new user without a home directory
RUN adduser -D -h /tmp/build/tests build

# Switch to the newly created user
USER build

# Copy the local package files to the container's workspace
COPY --chown=build . .

# Install any dependencies
RUN go get -d -v ./...

# Command to run the executable
CMD sh -c "go test -p 1 -v ./..."