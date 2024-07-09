# syntax=docker/dockerfile:1
# A sample microservice in Go packaged into a container image.
FROM golang:1.22-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY ./ ./
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /muzz-app

# Run
CMD ["/muzz-app"]