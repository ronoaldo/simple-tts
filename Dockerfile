# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.22-bookworm as builder

# Copy local code to the container image.
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the command inside the container.
# -o /app/server places the executable in the /app directory.
# CGO_ENABLED=0 builds a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/tts

# Use a slim base image to reduce the final image size.
# "gcr.io/distroless/base-debian11" is a good choice.
FROM gcr.io/distroless/base-debian11

# Copy the binary to the production image.
COPY --from=builder /app/server /server

# Set the entrypoint to run the binary.
ENTRYPOINT ["/server"]
