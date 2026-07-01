# ============================================
# 后端 Go 构建
# ============================================
FROM golang:1.22-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

# 先复制依赖文件
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || go mod tidy

# 复制源码
COPY . .
RUN go build -ldflags="-s -w" -o /app/pcms-server main.go

# ============================================
# 运行镜像
# ============================================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/pcms-server .
RUN mkdir -p uploads

EXPOSE 8080

CMD ["./pcms-server"]
