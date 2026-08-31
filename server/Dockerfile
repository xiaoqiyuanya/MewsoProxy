# 构建阶段
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -trimpath -ldflags="-s -w" -o /app/server .

# 运行阶段
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && mkdir -p /app/run \
    && chown -R 10001:10001 /app/run
WORKDIR /app
COPY --from=builder /app/server .
COPY config.yaml ./config.yaml
USER 10001:10001
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/healthz || exit 1
ENTRYPOINT ["/app/server"]
