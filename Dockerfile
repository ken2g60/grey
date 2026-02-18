FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed if using private/public modules)
RUN apk add --no-cache git

# Copy ONLY go.mod and go.sum first (for layer caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]
