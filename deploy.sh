#!/bin/bash
# ============================================
# PCMS 一键部署脚本 (Linux / macOS)
# ============================================

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info()  { echo -e "${GREEN}[INFO]${NC}  $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_step()  { echo -e "\n${BLUE}============================================${NC}"; echo -e "${BLUE}  $1${NC}"; echo -e "${BLUE}============================================${NC}\n"; }

# ============================================
# 检查必要工具
# ============================================
check_prerequisites() {
    log_step "检查环境依赖"

    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    log_info "Docker ✓"

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装"
        exit 1
    fi
    log_info "Docker Compose ✓"
}

# ============================================
# 生成 .env 文件
# ============================================
generate_env() {
    if [ ! -f ".env" ]; then
        log_info "生成 .env 配置文件..."
        cat > .env << EOF
# 数据库配置
DB_PORT=5342

# 后端配置
SERVER_PORT=8081

# 前端配置
WEB_PORT=80

# JWT 密钥 (生产环境请修改为复杂随机字符串)
JWT_SECRET=$(openssl rand -base64 32 2>/dev/null || echo "pcms-jwt-$(date +%s)")
EOF
        log_info ".env 文件已生成"
    else
        log_info ".env 文件已存在，跳过"
    fi
}

# ============================================
# 构建并启动所有服务
# ============================================
start_services() {
    log_step "构建并启动服务"

    # 拉取基础镜像
    log_info "拉取基础镜像..."
    docker pull postgres:16-alpine
    docker pull golang:1.25-alpine
    docker pull node:20-alpine
    docker pull alpine:3.19
    docker pull nginx:1.25-alpine

    # 构建并启动
    log_info "构建镜像..."
    docker compose build --no-cache

    log_info "启动容器..."
    docker compose up -d

    # 等待服务就绪
    log_info "等待服务启动..."
    sleep 5

    # 检查状态
    echo ""
    docker compose ps
    echo ""

    log_info "============================================"
    log_info "  PCMS 部署完成！"
    log_info "  前端地址: http://localhost:${WEB_PORT:-80}"
    log_info "  后端 API: http://localhost:${SERVER_PORT:-8081}/api/v1"
    log_info "============================================"
}

# ============================================
# 停止服务
# ============================================
stop_services() {
    log_step "停止所有服务"
    docker compose down
    log_info "服务已停止"
}

# ============================================
# 重启服务
# ============================================
restart_services() {
    log_step "重启所有服务"
    docker compose restart
    log_info "服务已重启"
}

# ============================================
# 查看日志
# ============================================
show_logs() {
    docker compose logs -f --tail=100 "$@"
}

# ============================================
# 重新构建（代码更新后）
# ============================================
rebuild() {
    log_step "重新构建并部署"
    docker compose down
    docker compose build --no-cache
    docker compose up -d
    log_info "重新构建完成"
}

# ============================================
# 清理
# ============================================
clean() {
    log_step "清理所有容器、镜像、数据卷"
    read -p "确认清理？这将删除所有数据！(yes/no): " confirm
    if [ "$confirm" = "yes" ]; then
        docker compose down -v --rmi all --remove-orphans
        log_info "清理完成"
    else
        log_info "已取消"
    fi
}

# ============================================
# 主入口
# ============================================
main() {
    case "${1:-deploy}" in
        deploy|start)
            check_prerequisites
            generate_env
            start_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        rebuild)
            rebuild
            ;;
        logs)
            shift
            show_logs "$@"
            ;;
        status)
            docker compose ps
            ;;
        clean)
            clean
            ;;
        *)
            echo "用法: $0 {deploy|start|stop|restart|rebuild|logs|status|clean}"
            echo ""
            echo "  deploy    - 一键部署（默认）"
            echo "  start     - 启动服务"
            echo "  stop      - 停止服务"
            echo "  restart   - 重启服务"
            echo "  rebuild   - 更新代码后重新构建"
            echo "  logs      - 查看日志 (可加容器名，如: $0 logs server)"
            echo "  status    - 查看服务状态"
            echo "  clean     - 清理所有（含数据）"
            exit 1
            ;;
    esac
}

main "$@"
