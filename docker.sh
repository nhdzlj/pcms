#!/bin/bash
# ============================================
# PCMS Docker 管理脚本 (Linux / macOS)
# 纯 Docker 操作，不拉取代码
# ============================================

set -e
cd "$(dirname "$0")"

GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC}  $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_step()  { echo -e "\n${BLUE}==> ${1}${NC}\n"; }

# ============================================
# 构建并启动
# ============================================
up() {
    log_step "构建并启动所有服务"
    docker compose build
    docker compose up -d

    echo ""
    log_info "等待服务就绪..."
    sleep 3
    docker compose ps
    echo ""
    log_info "部署完成！访问 http://localhost"
}

# ============================================
# 仅构建（不启动）
# ============================================
build() {
    log_step "仅构建镜像"
    docker compose build --no-cache
    log_info "构建完成"
}

# ============================================
# 启动（使用已有镜像）
# ============================================
start() {
    log_step "启动所有服务"
    docker compose up -d
    docker compose ps
    log_info "启动完成"
}

# ============================================
# 停止
# ============================================
stop() {
    log_step "停止所有服务"
    docker compose down
    log_info "已停止"
}

# ============================================
# 重启
# ============================================
restart() {
    log_step "重启所有服务"
    docker compose restart
    log_info "已重启"
}

# ============================================
# 查看日志
# ============================================
logs() {
    local svc="${1:-}"
    if [ -n "$svc" ]; then
        docker compose logs -f --tail=100 "$svc"
    else
        docker compose logs -f --tail=100
    fi
}

# ============================================
# 状态查看
# ============================================
status() {
    docker compose ps
}

# ============================================
# 进入容器
# ============================================
exec_cmd() {
    local svc="${1:-server}"
    docker compose exec "$svc" sh
}

# ============================================
# 主入口
# ============================================
case "${1:-up}" in
    up)       up ;;
    build)    build ;;
    start)    start ;;
    stop)     stop ;;
    restart)  restart ;;
    logs)     logs "$2" ;;
    status)   status ;;
    exec)     exec_cmd "$2" ;;
    *)
        echo "用法: $0 {up|build|start|stop|restart|logs|status|exec}"
        echo ""
        echo "  up       - 构建镜像并启动所有服务"
        echo "  build    - 仅构建镜像（不启动）"
        echo "  start    - 启动已有服务"
        echo "  stop     - 停止所有服务"
        echo "  restart  - 重启所有服务"
        echo "  logs     - 查看日志（可指定容器：$0 logs server）"
        echo "  status   - 查看服务状态"
        echo "  exec     - 进入容器 shell（默认 server）"
        ;;
esac
