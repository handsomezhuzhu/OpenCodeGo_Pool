# syntax=docker/dockerfile:1

# ---------- 阶段 1：构建前端 ----------
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend
WORKDIR /build/web

# 仅拷贝依赖清单，充分利用构建缓存
COPY web/package.json web/package-lock.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

# ---------- 阶段 2：编译 Go 二进制 ----------
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS backend
ARG TARGETOS TARGETARCH
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 前端产物嵌入 internal/frontend/dist（embed.FS）
COPY --from=frontend /build/web/dist ./internal/frontend/dist

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/opencode-pool ./cmd/server

# ---------- 阶段 3：运行时 ----------
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 1000 app

WORKDIR /app
COPY --from=backend /out/opencode-pool ./opencode-pool

# 数据库默认路径为 ./data/pool.db（相对工作目录）
RUN mkdir -p /app/data && chown -R app:app /app

USER app
EXPOSE 8080
VOLUME ["/app/data"]

# 将 config.yaml 挂载到 /app/config.yaml 即可自定义配置（只读即可）
ENTRYPOINT ["./opencode-pool"]
