# deploy/chat/development/Dockerfile
FROM golang:1.24.5

WORKDIR /app

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Expose port
EXPOSE 2525

# Run with air for hot reloading
CMD ["air", "-c", ".air/.air-chat.toml"]