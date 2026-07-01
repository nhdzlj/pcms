# ============================================
# 后端 Go 构建
# ============================================
FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY go.mod ./
COPY . .
RUN go mod tidy
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
