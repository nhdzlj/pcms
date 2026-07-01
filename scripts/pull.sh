#!/bin/bash
# ============================================
# PCMS GitHub 拉取 + 自动部署脚本 (Linux / macOS)
# ============================================

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC}  $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# ============================================
# 配置 (请修改为你的仓库地址)
# ============================================
GITHUB_REPO="${GITHUB_REPO:-https://github.com/YOUR_USERNAME/pcms.git}"
BRANCH="${BRANCH:-main}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# ============================================
# 使用方式提示
# ============================================
usage() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  pull        拉取最新代码（不重启）"
    echo "  deploy      拉取代码 + 重新构建部署"
    echo "  clone       首次克隆项目"
    echo ""
    echo "环境变量:"
    echo "  GITHUB_REPO  GitHub 仓库地址"
    echo "  BRANCH       分支名 (默认: main)"
    echo ""
    echo "示例:"
    echo "  $0 pull"
    echo "  $0 deploy"
    echo "  GITHUB_REPO=https://github.com/user/pcms.git $0 clone"
}

# ============================================
# 首次克隆
# ============================================
do_clone() {
    log_info "克隆项目: $GITHUB_REPO"
    git clone "$GITHUB_REPO" "$PROJECT_DIR" || {
        log_error "克隆失败，请检查仓库地址和权限"
        exit 1
    }
    log_info "克隆完成！"
}

# ============================================
# 拉取更新
# ============================================
do_pull() {
    log_info "切换到项目目录: $PROJECT_DIR"
    cd "$PROJECT_DIR"

    if [ ! -d ".git" ]; then
        log_error "不是 Git 仓库，请先使用 clone 命令克隆项目"
        log_error "或者设置 GITHUB_REPO 环境变量后运行: $0 clone"
        exit 1
    fi

    log_info "当前分支: $(git rev-parse --abbrev-ref HEAD)"
    log_info "拉取远程: $GITHUB_REPO (分支: $BRANCH)"

    # 保存本地修改
    git stash save "auto-stash-before-pull-$(date +%Y%m%d%H%M%S)" 2>/dev/null || true

    # 拉取最新代码
    git fetch origin "$BRANCH"
    git checkout "$BRANCH" 2>/dev/null || true
    git reset --hard "origin/$BRANCH"

    log_info "代码拉取完成"
    git log --oneline -5
}

# ============================================
# 主入口
# ============================================
case "${1:-pull}" in
    clone)
        do_clone
        ;;
    pull)
        do_pull
        ;;
    deploy)
        do_pull
        log_info "正在重新构建部署..."
        cd "$PROJECT_DIR"
        bash deploy.sh rebuild
        ;;
    *)
        usage
        exit 1
        ;;
esac
