FROM golang:1.22-alpine

RUN apk add --no-cache git make

WORKDIR /app

# Copy source
COPY . .

# Download dependencies
RUN go mod tidy && go mod download

# Default command: build and test
CMD ["sh", "-c", "go build ./... && go test -v -race ./... && go vet ./... && echo '==> All checks passed!'"]
