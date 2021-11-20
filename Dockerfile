# syntax = docker/dockerfile:1.2.1

# Start from the latest golang base image
FROM golang:1.17.3-buster as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN make build

FROM bitnami/minideb:buster

RUN install_packages ca-certificates && \
    update-ca-certificates

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/ .

# ensure that we'll not be running the container as root
USER 1001:1001

# Command to run the executable
ENTRYPOINT ["./httplogger"]
