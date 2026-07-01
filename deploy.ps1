# ============================================
# PCMS 一键部署脚本 (Windows PowerShell)
# ============================================

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

# 颜色定义
function Write-Info  { Write-Host "[INFO]  " -ForegroundColor Green -NoNewline; Write-Host $args[0] }
function Write-Warn  { Write-Host "[WARN]  " -ForegroundColor Yellow -NoNewline; Write-Host $args[0] }
function Write-Error { Write-Host "[ERROR] " -ForegroundColor Red   -NoNewline; Write-Host $args[0] }
function Write-Step  {
    Write-Host ""
    Write-Host "============================================" -ForegroundColor Blue
    Write-Host "  $args[0]" -ForegroundColor Blue
    Write-Host "============================================" -ForegroundColor Blue
    Write-Host ""
}

# ============================================
# 检查必要工具
# ============================================
function Check-Prerequisites {
    Write-Step "检查环境依赖"

    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-Error "Docker 未安装，请先安装 Docker Desktop"
        exit 1
    }
    Write-Info "Docker ✓"

    # Docker Compose v2 (docker compose)
    $composeCmd = $null
    if (Get-Command docker -ErrorAction SilentlyContinue) {
        docker compose version 2>$null | Out-Null
        if ($LASTEXITCODE -eq 0) {
            $composeCmd = "docker compose"
        }
    }
    if (-not $composeCmd) {
        Write-Error "Docker Compose 未找到"
        exit 1
    }
    Write-Info "Docker Compose ✓"
}

# ============================================
# 生成 .env 文件
# ============================================
function Generate-Env {
    if (-not (Test-Path ".env")) {
        Write-Info "生成 .env 配置文件..."
        $randomSecret = -join ((48..57) + (65..90) + (97..122) | Get-Random -Count 32 | ForEach-Object { [char]$_ })
        @"
# 数据库配置
DB_PORT=5432

# 后端配置
SERVER_PORT=8081

# 前端配置
WEB_PORT=80

# JWT 密钥 (生产环境请修改为复杂随机字符串)
JWT_SECRET=$randomSecret
"@ | Out-File -FilePath ".env" -Encoding UTF8
        Write-Info ".env 文件已生成"
    } else {
        Write-Info ".env 文件已存在，跳过"
    }
}

# ============================================
# 构建并启动所有服务
# ============================================
function Start-Services {
    Write-Step "构建并启动服务"

    Write-Info "拉取基础镜像..."
    docker pull postgres:16-alpine
    docker pull golang:1.22-alpine
    docker pull node:20-alpine
    docker pull alpine:3.19
    docker pull nginx:1.25-alpine

    Write-Info "构建镜像..."
    docker compose build --no-cache

    Write-Info "启动容器..."
    docker compose up -d

    Write-Info "等待服务启动..."
    Start-Sleep -Seconds 5

    Write-Host ""
    docker compose ps
    Write-Host ""

    $webPort = if ($env:WEB_PORT) { $env:WEB_PORT } else { "80" }
    $serverPort = if ($env:SERVER_PORT) { $env:SERVER_PORT } else { "8081" }

    Write-Info "============================================"
    Write-Info "  PCMS 部署完成！"
    Write-Info "  前端地址: http://localhost:$webPort"
    Write-Info "  后端 API: http://localhost:$serverPort/api/v1"
    Write-Info "============================================"
}

# ============================================
# 停止服务
# ============================================
function Stop-Services {
    Write-Step "停止所有服务"
    docker compose down
    Write-Info "服务已停止"
}

# ============================================
# 重启服务
# ============================================
function Restart-Services {
    Write-Step "重启所有服务"
    docker compose restart
    Write-Info "服务已重启"
}

# ============================================
# 查看日志
# ============================================
function Show-Logs {
    param([string]$Container)
    if ($Container) {
        docker compose logs -f --tail=100 $Container
    } else {
        docker compose logs -f --tail=100
    }
}

# ============================================
# 重新构建（代码更新后）
# ============================================
function Rebuild {
    Write-Step "重新构建并部署"
    docker compose down
    docker compose build --no-cache
    docker compose up -d
    Write-Info "重新构建完成"
}

# ============================================
# 清理
# ============================================
function Clean {
    Write-Step "清理所有容器、镜像、数据卷"
    $confirm = Read-Host "确认清理？这将删除所有数据！(yes/no)"
    if ($confirm -eq "yes") {
        docker compose down -v --rmi all --remove-orphans
        Write-Info "清理完成"
    } else {
        Write-Info "已取消"
    }
}

# ============================================
# 主入口
# ============================================
function Main {
    $action = if ($args.Count -gt 0) { $args[0] } else { "deploy" }
    $container = if ($args.Count -gt 1) { $args[1] } else { $null }

    switch ($action) {
        "deploy"  { Check-Prerequisites; Generate-Env; Start-Services }
        "start"   { Check-Prerequisites; Start-Services }
        "stop"    { Stop-Services }
        "restart" { Restart-Services }
        "rebuild" { Rebuild }
        "logs"    { Show-Logs -Container $container }
        "status"  { docker compose ps }
        "clean"   { Clean }
        default {
            Write-Host "用法: .\deploy.ps1 {deploy|start|stop|restart|rebuild|logs|status|clean}"
            Write-Host ""
            Write-Host "  deploy   - 一键部署（默认）"
            Write-Host "  start    - 启动服务"
            Write-Host "  stop     - 停止服务"
            Write-Host "  restart  - 重启服务"
            Write-Host "  rebuild  - 更新代码后重新构建"
            Write-Host "  logs     - 查看日志 (可加容器名，如: .\deploy.ps1 logs server)"
            Write-Host "  status   - 查看服务状态"
            Write-Host "  clean    - 清理所有（含数据）"
        }
    }
}

Main @args
