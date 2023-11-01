# Use the official Golang image to create a build artifact.
FROM golang:1.18 AS builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Build the binary.
RUN go build -v -o server

# Use a Docker multi-stage build to create a lean production image.
FROM golang:1.18

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Run the web service on container startup.
CMD ["/server"]
