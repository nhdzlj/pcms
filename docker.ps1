# ============================================
# PCMS Docker 管理脚本 (Windows PowerShell)
# 纯 Docker 操作，不拉取代码
# ============================================

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

function Write-Info  { Write-Host "[INFO]  " -ForegroundColor Green -NoNewline; Write-Host $args[0] }
function Write-Error { Write-Host "[ERROR] " -ForegroundColor Red   -NoNewline; Write-Host $args[0] }
function Write-Step  {
    Write-Host "`n==> $args[0]`n" -ForegroundColor Blue
}

# ============================================
# 构建并启动
# ============================================
function Invoke-Up {
    Write-Step "构建并启动所有服务"
    docker compose build
    docker compose up -d

    Write-Host ""
    Write-Info "等待服务就绪..."
    Start-Sleep -Seconds 3
    docker compose ps
    Write-Host ""
    Write-Info "部署完成！访问 http://localhost"
}

# ============================================
# 仅构建（不启动）
# ============================================
function Invoke-Build {
    Write-Step "仅构建镜像"
    docker compose build --no-cache
    Write-Info "构建完成"
}

# ============================================
# 启动（使用已有镜像）
# ============================================
function Invoke-Start {
    Write-Step "启动所有服务"
    docker compose up -d
    docker compose ps
    Write-Info "启动完成"
}

# ============================================
# 停止
# ============================================
function Invoke-Stop {
    Write-Step "停止所有服务"
    docker compose down
    Write-Info "已停止"
}

# ============================================
# 重启
# ============================================
function Invoke-Restart {
    Write-Step "重启所有服务"
    docker compose restart
    Write-Info "已重启"
}

# ============================================
# 查看日志
# ============================================
function Invoke-Logs {
    param([string]$Service)
    if ($Service) {
        docker compose logs -f --tail=100 $Service
    } else {
        docker compose logs -f --tail=100
    }
}

# ============================================
# 状态查看
# ============================================
function Invoke-Status {
    docker compose ps
}

# ============================================
# 进入容器
# ============================================
function Invoke-Exec {
    param([string]$Service = "server")
    docker compose exec $Service sh
}

# ============================================
# 主入口
# ============================================
$action = if ($args.Count -gt 0) { $args[0] } else { "up" }
$arg2   = if ($args.Count -gt 1) { $args[1] } else { $null }

switch ($action) {
    "up"      { Invoke-Up }
    "build"   { Invoke-Build }
    "start"   { Invoke-Start }
    "stop"    { Invoke-Stop }
    "restart" { Invoke-Restart }
    "logs"    { Invoke-Logs -Service $arg2 }
    "status"  { Invoke-Status }
    "exec"    { Invoke-Exec -Service $arg2 }
    default {
        Write-Host "用法: .\docker.ps1 {up|build|start|stop|restart|logs|status|exec}"
        Write-Host ""
        Write-Host "  up       - 构建镜像并启动所有服务"
        Write-Host "  build    - 仅构建镜像（不启动）"
        Write-Host "  start    - 启动已有服务"
        Write-Host "  stop     - 停止所有服务"
        Write-Host "  restart  - 重启所有服务"
        Write-Host "  logs     - 查看日志（可指定容器：.\docker.ps1 logs server）"
        Write-Host "  status   - 查看服务状态"
        Write-Host "  exec     - 进入容器 shell（默认 server）"
    }
}
