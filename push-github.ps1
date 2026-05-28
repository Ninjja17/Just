# Push to GitHub Script
# Repository: https://github.com/Ninjja17/Just.git

Write-Host "Starting Git Push Process..." -ForegroundColor Cyan
Write-Host ""

# Check Git
Write-Host "[1] Checking Git..." -ForegroundColor Yellow
try {
    $version = git --version
    Write-Host "    OK - $version" -ForegroundColor Green
} catch {
    Write-Host "    ERROR - Git not found!" -ForegroundColor Red
    Write-Host "    Install from: https://git-scm.com/download/win" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit
}

Write-Host ""

# Configure Git
Write-Host "[2] Configuring Git..." -ForegroundColor Yellow
$userName = git config --global user.name
$userEmail = git config --global user.email

if (-not $userName -or -not $userEmail) {
    Write-Host "    Git not configured yet" -ForegroundColor Yellow
    $name = Read-Host "    Enter your name"
    $email = Read-Host "    Enter your email"
    git config --global user.name "$name"
    git config --global user.email "$email"
    Write-Host "    OK - Git configured" -ForegroundColor Green
} else {
    Write-Host "    OK - Already configured" -ForegroundColor Green
}

Write-Host ""

# Initialize repo
Write-Host "[3] Initializing repository..." -ForegroundColor Yellow
if (Test-Path ".git") {
    Write-Host "    OK - Repository exists" -ForegroundColor Gray
} else {
    git init
    Write-Host "    OK - Repository created" -ForegroundColor Green
}

Write-Host ""

# Add files
Write-Host "[4] Adding files..." -ForegroundColor Yellow
git add .
Write-Host "    OK - All files staged" -ForegroundColor Green

Write-Host ""

# Commit
Write-Host "[5] Creating commit..." -ForegroundColor Yellow
git commit -m "Initial commit: 10000 Hour Tracker application"
Write-Host "    OK - Commit created" -ForegroundColor Green

Write-Host ""

# Add remote
Write-Host "[6] Adding remote..." -ForegroundColor Yellow
$remoteExists = git remote get-url origin 2>$null
if ($remoteExists) {
    Write-Host "    INFO - Remote exists" -ForegroundColor Gray
} else {
    git remote add origin https://github.com/Ninjja17/Just.git
    Write-Host "    OK - Remote added" -ForegroundColor Green
}

Write-Host ""

# Set branch
Write-Host "[7] Setting main branch..." -ForegroundColor Yellow
git branch -M main
Write-Host "    OK - Branch set to main" -ForegroundColor Green

Write-Host ""

# Push
Write-Host "[8] Pushing to GitHub..." -ForegroundColor Yellow
Write-Host "    URL: https://github.com/Ninjja17/Just.git" -ForegroundColor Cyan
Write-Host ""
Write-Host "    Authentication required:" -ForegroundColor Yellow
Write-Host "    Username: Your GitHub username" -ForegroundColor Gray
Write-Host "    Password: Personal Access Token" -ForegroundColor Gray
Write-Host "    Get token: https://github.com/settings/tokens" -ForegroundColor Gray
Write-Host ""

try {
    git push -u origin main
    Write-Host ""
    Write-Host "    SUCCESS!" -ForegroundColor Green
    Write-Host ""
    Write-Host "    Your project is now on GitHub!" -ForegroundColor Green
    Write-Host "    Visit: https://github.com/Ninjja17/Just" -ForegroundColor Cyan
} catch {
    Write-Host ""
    Write-Host "    ERROR - Push failed!" -ForegroundColor Red
    Write-Host "    Try:" -ForegroundColor Yellow
    Write-Host "    1. Use Personal Access Token as password" -ForegroundColor Gray
    Write-Host "    2. Check repository exists on GitHub" -ForegroundColor Gray
}

Write-Host ""
Read-Host "Press Enter to exit"
