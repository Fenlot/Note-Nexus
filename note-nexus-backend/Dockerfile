# ============================
# STAGE 1: Build the Go Binary
# ============================
# Use a solid, recent Go version on Alpine Linux (small size)
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container build area
WORKDIR /app

# Copy dependency files first. Docker caches these layers, 
# so if you only change code, it won't re-download modules.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application.
# CGO_ENABLED=0 creates a statically linked binary (no external dependencies).
# This is possible because we used the 'modernc.org/sqlite' pure Go driver!
RUN CGO_ENABLED=0 GOOS=linux go build -o note-nexus-api ./cmd/api

# ============================
# STAGE 2: The Runtime Image
# ============================
# Start with a tiny, empty Alpine Linux image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from Stage 1
COPY --from=builder /app/note-nexus-api .

# IMPORTANT: Copy the migrations folder so the app can set up the DB
COPY --from=builder /app/migrations ./migrations

# Tell Docker this container listens on port 8080
EXPOSE 8080

# The command to run when the container starts
CMD ["./note-nexus-api"]