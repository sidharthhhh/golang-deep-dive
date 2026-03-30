# ============================================================
# Stage 1: Build
# ============================================================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /app/auth-service \
    cmd/server/main.go

# ============================================================
# Stage 2: Runtime (minimal image, no nginx)
# ============================================================
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata wget

RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

WORKDIR /home/appuser

COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appgroup /home/appuser

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./auth-service"]
