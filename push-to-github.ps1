# ====================================
# Automated Git Push Script
# Repository: https://github.com/Ninjja17/Just.git
# ====================================

Write-Host "Starting Git Push Process..." -ForegroundColor Cyan
Write-Host ""

# Check if Git is installed
Write-Host "[1] Checking Git installation..." -ForegroundColor Yellow
try {
    $gitVersion = git --version
    Write-Host "   [OK] Git found: $gitVersion" -ForegroundColor Green
} catch {
    Write-Host "   [ERROR] Git is not installed!" -ForegroundColor Red
    Write-Host "   Please install Git from: https://git-scm.com/download/win" -ForegroundColor Yellow
    Write-Host "   Then restart PowerShell and run this script again." -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit
}

Write-Host ""

# Check if Git is configured
Write-Host "[2] Checking Git configuration..." -ForegroundColor Yellow
$userName = git config --global user.name
$userEmail = git config --global user.email

if (-not $userName -or -not $userEmail) {
    Write-Host "   [WARNING] Git not configured. Let's set it up!" -ForegroundColor Yellow
    Write-Host ""
    
    $name = Read-Host "   Enter your name"
    $email = Read-Host "   Enter your email"
    
    git config --global user.name "$name"
    git config --global user.email "$email"
    
    Write-Host "   [OK] Git configured successfully!" -ForegroundColor Green
} else {
    Write-Host "   [OK] Git already configured" -ForegroundColor Green
    Write-Host "   Name: $userName" -ForegroundColor Gray
    Write-Host "   Email: $userEmail" -ForegroundColor Gray
}

Write-Host ""

# Initialize Git repository
Write-Host "[3] Initializing Git repository..." -ForegroundColor Yellow
if (Test-Path ".git") {
    Write-Host "   [INFO] Git repository already exists" -ForegroundColor Gray
} else {
    git init
    Write-Host "   [OK] Repository initialized" -ForegroundColor Green
}

Write-Host ""

# Add all files
Write-Host "[4] Adding files to Git..." -ForegroundColor Yellow
git add .
$fileCount = (git diff --cached --numstat | Measure-Object).Count
Write-Host "   [OK] Added all files (staging $fileCount changes)" -ForegroundColor Green

Write-Host ""

# Create commit
Write-Host "[5] Creating commit..." -ForegroundColor Yellow
$commitMessage = "Initial commit: 10000 Hour Tracker - Full-stack application

Features:
- Go backend with Gin framework
- React TypeScript frontend with Tailwind CSS
- PostgreSQL database with complete schema
- JWT authentication with OTP verification
- Google OAuth integration (ready)
- Redis caching for sessions
- Docker Compose setup
- Complete API with 30+ endpoints

Modules:
- Authentication & User Management
- Skills tracking
- Practice sessions with timer
- Goals & milestones
- Analytics & insights
- Social features (follow, leaderboard)
- Notifications system

Tech Stack:
Backend: Go, Gin, PostgreSQL, Redis, JWT, SendGrid
Frontend: React 18, TypeScript, Vite, TailwindCSS, Axios
DevOps: Docker, Docker Compose"

git commit -m "$commitMessage"
Write-Host "   [OK] Commit created successfully" -ForegroundColor Green

Write-Host ""

# Add remote origin
Write-Host "[6] Adding remote repository..." -ForegroundColor Yellow
$remoteExists = git remote get-url origin 2>$null
if ($remoteExists) {
    Write-Host "   ℹ️  Remote 'origin' already exists: $remoteExists" -ForegroundColor Gray
    $response = Read-Host "   Do you want to update it? (y/n)"
    if ($response -[INFO] Remote 'origin' already exists: $remoteExists" -ForegroundColor Gray
    $response = Read-Host "   Do you want to update it? (y/n)"
    if ($response -eq "y" -or $response -eq "Y") {
        git remote set-url origin https://github.com/Ninjja17/Just.git
        Write-Host "   [OK] Remote updated" -ForegroundColor Green
    }
} else {
    git remote add origin https://github.com/Ninjja17/Just.git
    Write-Host "   [OK] Remote added: https://github.com/Ninjja17/Just.git" -ForegroundColor Green
}

Write-Host ""

# Rename branch to main
Write-Host "[7] Setting up main branch..." -ForegroundColor Yellow
git branch -M main
Write-Host "   [OK] Branch set to 'main'" -ForegroundColor Green

Write-Host ""

# Push to GitHub
Write-Host "[8] Pushing to GitHub..." -ForegroundColor Yellow
Write-Host "   Pushing to: https://github.com/Ninjja17/Just.git" -ForegroundColor Cyan
Write-Host ""
Write-Host "   [WARNING] You may be prompted for authentication:" -ForegroundColor Yellow
Write-Host "   - Username: Your GitHub username" -ForegroundColor Gray
Write-Host "   - Password: Your Personal Access Token (NOT your GitHub password)" -ForegroundColor Gray
Write-Host ""
Write-Host "   TIP:

try {
    git push -u origin main
    Write-Host ""
    Write-Host "   [OK] Successfully pushed to GitHub!" -ForegroundColor Green
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host " SUCCESS! Your project is now on GitHub!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Repository URL: https://github.com/Ninjja17/Just" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Project contains:" -ForegroundColor White
    Write-Host "   - Go backend" -ForegroundColor Gray
    Write-Host "   - React frontend" -ForegroundColor Gray
    Write-Host "   - PostgreSQL migrations" -ForegroundColor Gray
    Write-Host "   - Docker configuration" -ForegroundColor Gray
    Write-Host "   - Complete documentation" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host ""
    Write-Host "   [ERROR] Push failed!" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "   Common solutions:" -ForegroundColor Yellow
    Write-Host "   1. Make sure you have push access to the repository" -ForegroundColor Gray
    Write-Host "   2. Use a Personal Access Token instead of password" -ForegroundColor Gray
    Write-Host "   3. Check if the repository exists on GitHub" -ForegroundColor Gray
    Write-Host ""
}

Write-Host ""
Write-Host "Press Enter to exit" -ForegroundColor Gray
Read-Host
