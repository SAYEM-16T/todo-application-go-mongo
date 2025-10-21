# syntax=docker/dockerfile:1.6

########################
# BUILD STAGE
########################
FROM golang:1.24 AS builder
WORKDIR /app/backend

# Copy backend + frontend
COPY backend/ /app/backend/
COPY frontend/ /app/frontend/

# Prevent VCS warning
ENV GOFLAGS=-buildvcs=false

# Download dependencies
RUN go mod tidy

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

########################
# RUNTIME STAGE
########################
FROM alpine:3.20

# Create non-root user
RUN adduser -D -u 10001 app
USER app

WORKDIR /app/backend

# Copy server binary + static assets
COPY --from=builder /app/server /app/server
COPY --from=builder /app/frontend /app/frontend

EXPOSE 8080

# Default environment
ENV ADDR=:8080 \
    DB_NAME=todoapp

ENTRYPOINT ["/app/server"]
