# 🚀 PUSH TO GITHUB - QUICK START

Your repository: **https://github.com/Ninjja17/Just.git**

## ⚡ ONE-COMMAND PUSH

### Step 1: Install Git (if not installed)
Download and install: **https://git-scm.com/download/win**
- Use default settings
- Restart PowerShell after installation

### Step 2: Run the Script
Open PowerShell in this directory and run:

```powershell
.\push-to-github.ps1
```

That's it! The script will:
- ✅ Check Git installation
- ✅ Configure Git (if needed)
- ✅ Initialize repository
- ✅ Add all files (65 files)
- ✅ Create commit with detailed message
- ✅ Connect to your GitHub repo
- ✅ Push everything to GitHub

## 🔐 Authentication

When prompted, use:
- **Username:** `Ninjja17` (your GitHub username)
- **Password:** Your **Personal Access Token** (NOT your GitHub password)

### Get a Personal Access Token:
1. Go to: https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Give it a name (e.g., "10000hr-tracker")
4. Check the `repo` permission box
5. Click "Generate token"
6. Copy the token and use it as your password

## 🎯 Alternative: Manual Commands

If you prefer to run commands manually:

```powershell
# 1. Initialize
git init

# 2. Configure (first time only)
git config --global user.name "Your Name"
git config --global user.email "your@email.com"

# 3. Add files
git add .

# 4. Commit
git commit -m "Initial commit: 10000 Hour Tracker application"

# 5. Add remote
git remote add origin https://github.com/Ninjja17/Just.git

# 6. Push
git branch -M main
git push -u origin main
```

## ✅ After Successful Push

Visit your repository:
**https://github.com/Ninjja17/Just**

You should see all 65 files uploaded! 🎉

## 🆘 Troubleshooting

**"Git not found"**
→ Install Git and restart PowerShell

**"Permission denied"**
→ Make sure you're using a Personal Access Token, not your GitHub password

**"Repository not found"**
→ Check if the repository exists on GitHub
→ Make sure you have access to it

**"Script execution disabled"**
→ Run: `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser`

---

**Need help?** Check these files:
- [GITHUB_GUIDE.md](GITHUB_GUIDE.md) - Complete guide
- [GITHUB_COMMANDS.md](GITHUB_COMMANDS.md) - All commands explained
