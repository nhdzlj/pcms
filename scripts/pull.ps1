# ============================================
# PCMS GitHub 拉取 + 自动部署脚本 (Windows PowerShell)
# ============================================

$ErrorActionPreference = "Stop"

function Write-Info  { Write-Host "[INFO]  " -ForegroundColor Green -NoNewline; Write-Host $args[0] }
function Write-Error { Write-Host "[ERROR] " -ForegroundColor Red   -NoNewline; Write-Host $args[0] }

# ============================================
# 配置 (请修改为你的仓库地址)
# ============================================
$GITHUB_REPO = if ($env:GITHUB_REPO) { $env:GITHUB_REPO } else { "https://github.com/YOUR_USERNAME/pcms.git" }
$BRANCH      = if ($env:BRANCH)      { $env:BRANCH }      else { "main" }

$ScriptDir  = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectDir = Split-Path -Parent $ScriptDir

# ============================================
# 使用方式提示
# ============================================
function Show-Usage {
    Write-Host "用法: .\scripts\pull.ps1 [选项]"
    Write-Host ""
    Write-Host "选项:"
    Write-Host "  pull        拉取最新代码（不重启）"
    Write-Host "  deploy      拉取代码 + 重新构建部署"
    Write-Host "  clone       首次克隆项目"
    Write-Host ""
    Write-Host "环境变量:"
    Write-Host "  `$env:GITHUB_REPO  GitHub 仓库地址"
    Write-Host "  `$env:BRANCH       分支名 (默认: main)"
    Write-Host ""
    Write-Host "示例:"
    Write-Host "  .\scripts\pull.ps1 pull"
    Write-Host "  `$env:GITHUB_REPO='https://github.com/user/pcms.git'; .\scripts\pull.ps1 clone"
}

# ============================================
# 首次克隆
# ============================================
function Invoke-Clone {
    Write-Info "克隆项目: $GITHUB_REPO"

    # 克隆到父目录（因为 scripts 目录已经在 pcms 内）
    $cloneDir = Split-Path -Parent $ProjectDir
    Set-Location $cloneDir

    git clone $GITHUB_REPO $ProjectDir
    if ($LASTEXITCODE -ne 0) {
        Write-Error "克隆失败，请检查仓库地址和权限"
        exit 1
    }
    Write-Info "克隆完成！"
}

# ============================================
# 拉取更新
# ============================================
function Invoke-Pull {
    Write-Info "切换到项目目录: $ProjectDir"
    Set-Location $ProjectDir

    if (-not (Test-Path ".git")) {
        Write-Error "不是 Git 仓库，请先使用 clone 命令克隆项目"
        exit 1
    }

    $currentBranch = git rev-parse --abbrev-ref HEAD
    Write-Info "当前分支: $currentBranch"
    Write-Info "拉取远程: $GITHUB_REPO (分支: $BRANCH)"

    # 保存本地修改
    $stashName = "auto-stash-before-pull-$(Get-Date -Format 'yyyyMMddHHmmss')"
    git stash save $stashName 2>$null

    # 拉取最新代码
    git fetch origin $BRANCH
    git checkout $BRANCH 2>$null
    git reset --hard "origin/$BRANCH"

    Write-Info "代码拉取完成"
    git log --oneline -5
}

# ============================================
# 主入口
# ============================================
$action = if ($args.Count -gt 0) { $args[0] } else { "pull" }

switch ($action) {
    "clone" {
        Invoke-Clone
    }
    "pull" {
        Invoke-Pull
    }
    "deploy" {
        Invoke-Pull
        Write-Info "正在重新构建部署..."
        Set-Location $ProjectDir
        & "$ProjectDir\deploy.ps1" rebuild
    }
    default {
        Show-Usage
        exit 1
    }
}
