FROM golang:alpine AS build

# Install build dependencies for CGo and UPX
RUN apk add --no-cache gcc musl-dev sqlite-dev upx

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

# Build with CGo enabled for SQLite support
RUN CGO_ENABLED=1 GOOS=linux go build -o backend -a -ldflags="-s -w -linkmode external -extldflags '-static'" .

RUN upx --ultra-brute -qq backend && upx -t backend

FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite-libs

# Create app directory and data directory
WORKDIR /app
RUN mkdir -p /app/data && chmod 755 /app/data

COPY --from=build /app/backend /app/backend
