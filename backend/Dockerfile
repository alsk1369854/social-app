# Frontend build stage
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend package files
COPY frontend/package*.json ./

# Install frontend dependencies
RUN npm ci

# Copy frontend source code
COPY frontend ./

# Build frontend
RUN npm run build

# Backend build stage
FROM golang:tip-alpine3.22 AS backend-builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend ./

# Generate Swagger docs (provides package `backend/docs`)
# 1) Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest
# 2) Generate docs at module root (/app)
RUN $(go env GOPATH)/bin/swag init -g main.go -o docs

# Build the application for Linux AMD64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main-linux-amd64 .

# Final stage
FROM golang:tip-alpine3.22

WORKDIR /root/

# Copy the binary from backend builder
COPY --from=backend-builder /app/main-linux-amd64 .

# Copy frontend build files to public directory
COPY --from=frontend-builder /app/frontend/build ./public

# Expose port (adjust if your Gin server uses a different port)
EXPOSE 8080

# Run the binary
CMD ["./main-linux-amd64", "--host=0.0.0.0", "--port=8080"]