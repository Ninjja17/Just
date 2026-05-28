# Quick Commands - Push to GitHub

## Prerequisites
✅ Git must be installed: https://git-scm.com/download/win

## 🚀 Commands to Run (In Order)

### 1️⃣ Open PowerShell in Project Directory
```powershell
cd C:\Users\SH40176270\Desktop\1000hr
```

### 2️⃣ Configure Git (First Time Only)
```powershell
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### 3️⃣ Initialize Git Repository
```powershell
git init
```

### 4️⃣ Check What Will Be Committed
```powershell
git status
```

### 5️⃣ Add All Files
```powershell
git add .
```

### 6️⃣ Create Initial Commit
```powershell
git commit -m "Initial commit: 10000 Hour Tracker full-stack application

- Go backend with JWT authentication and OTP verification
- React TypeScript frontend with Tailwind CSS
- PostgreSQL database schema with migrations
- Docker Compose setup for easy deployment
- Complete API with 30+ endpoints
- Authentication, Skills, Sessions, Goals, Analytics, Social features"
```

### 7️⃣ Create GitHub Repository
1. Go to: https://github.com/new
2. Repository name: `10000hr-tracker` (or your choice)
3. Description: `Track your journey to mastery - Full-stack time tracking app`
4. Choose Public or Private
5. **DO NOT** check any boxes (no README, .gitignore, or license)
6. Click "Create repository"

### 8️⃣ Connect to GitHub (Replace with YOUR info)
```powershell
# Replace YOUR_USERNAME and REPO_NAME with your actual GitHub username and repository name
git remote add origin https://github.com/YOUR_USERNAME/REPO_NAME.git

# Example:
# git remote add origin https://github.com/johndoe/10000hr-tracker.git
```

### 9️⃣ Rename Branch to Main
```powershell
git branch -M main
```

### 🔟 Push to GitHub
```powershell
git push -u origin main
```

## ✅ Done!

Visit your repository at: `https://github.com/YOUR_USERNAME/REPO_NAME`

---

## 🔄 For Future Updates

```powershell
# Add changes
git add .

# Commit with message
git commit -m "Description of changes"

# Push to GitHub
git push
```

## 🐛 Troubleshooting

**Authentication Required?**
- Use Personal Access Token instead of password
- Get token from: GitHub Settings → Developer settings → Personal access tokens

**Permission Denied?**
- Check if repository URL is correct
- Verify you have push access to the repository

**Already initialized?**
```powershell
# Check current remote
git remote -v

# Change remote if needed
git remote set-url origin https://github.com/YOUR_USERNAME/NEW_REPO.git
```

---

**Need detailed help? See [GITHUB_GUIDE.md](GITHUB_GUIDE.md)**
