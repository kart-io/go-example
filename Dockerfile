# Build stage
FROM golang:1.25-alpine AS builder

# Install git and ca-certificates (needed for go modules)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version injection
ARG SERVICE_NAME=go-example-api
ARG VERSION=unknown
ARG COMMIT=unknown
ARG BRANCH=unknown
ARG BUILD_DATE=unknown

# Build the application with version injection
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "\
    -w -s \
    -X 'github.com/kart-io/version.serviceName=${SERVICE_NAME}' \
    -X 'github.com/kart-io/version.gitVersion=${VERSION}' \
    -X 'github.com/kart-io/version.gitCommit=${COMMIT}' \
    -X 'github.com/kart-io/version.gitBranch=${BRANCH}' \
    -X 'github.com/kart-io/version.buildDate=${BUILD_DATE}' \
    " -o gin-demo ./gin-demo

# Runtime stage
FROM alpine:latest

# Install ca-certificates for TLS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gin-demo .

# Add labels for better container management
LABEL version="${VERSION}"
LABEL commit="${COMMIT}"  
LABEL branch="${BRANCH}"
LABEL build-date="${BUILD_DATE}"
LABEL service="${SERVICE_NAME}"

# Expose port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run the binary
CMD ["./gin-demo"]