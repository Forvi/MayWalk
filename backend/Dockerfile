# === Build ===
FROM golang:1.27.1-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
COPY api/openapi.yml ./docs/openapi.yml
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
RUN oapi-codegen -config oapi-codegen.yml docs/openapi.yml

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/swagger ./cmd/docs/main.go

# === Runtime ===
FROM alpine:3.19

RUN apk add --no-cache ca-certificates curl

RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -D appuser

WORKDIR /app

COPY --from=builder /app/bin/server ./server
COPY --from=builder /app/bin/swagger ./swagger

COPY backend/migrations ./migrations

COPY --from=builder /app/docs/static ./docs/static
COPY --from=builder /app/docs/openapi.yml ./docs/openapi.yml

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080 8081

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./server"]
