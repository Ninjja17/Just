# Push to GitHub Guide

## Step 1: Install Git

1. Download Git for Windows from: https://git-scm.com/download/win
2. Run the installer with default settings
3. Restart your terminal/PowerShell after installation

## Step 2: Configure Git (First Time Only)

Open PowerShell and run:

```powershell
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

## Step 3: Initialize Git Repository

Navigate to your project and initialize git:

```powershell
cd C:\Users\SH40176270\Desktop\1000hr
git init
```

## Step 4: Add Files to Git

```powershell
# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: 10000 Hour Tracker - Full stack application"
```

## Step 5: Create GitHub Repository

1. Go to GitHub: https://github.com
2. Click the **"+"** icon in the top right
3. Select **"New repository"**
4. Fill in details:
   - **Repository name:** `10000hr-tracker` (or your preferred name)
   - **Description:** `A web application for tracking your journey to mastery through the 10,000-hour rule`
   - **Visibility:** Choose Public or Private
   - **❌ DO NOT** initialize with README, .gitignore, or license (we already have these)
5. Click **"Create repository"**

## Step 6: Connect to GitHub and Push

After creating the repository on GitHub, you'll see a page with commands. Use these:

```powershell
# Add the remote repository (replace YOUR_USERNAME and YOUR_REPO with your actual values)
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO.git

# Rename branch to main (if needed)
git branch -M main

# Push to GitHub
git push -u origin main
```

### Example:
If your GitHub username is `johndoe` and repo name is `10000hr-tracker`:

```powershell
git remote add origin https://github.com/johndoe/10000hr-tracker.git
git branch -M main
git push -u origin main
```

## Step 7: Verify Upload

Visit your GitHub repository URL to see all your files!

## 📝 What Gets Uploaded:

✅ **Backend** (Go files)
✅ **Frontend** (React/TypeScript files)
✅ **Docker** configuration
✅ **Documentation** (README, SETUP, etc.)
✅ **Migrations** (Database schema)

❌ **Not uploaded** (thanks to .gitignore):
- node_modules/
- .env files (environment secrets)
- Build outputs
- Database files
- Temporary files

## 🔄 Future Updates

After making changes, push updates with:

```powershell
git add .
git commit -m "Description of your changes"
git push
```

## 🌿 Working with Branches (Optional)

Create feature branches for new work:

```powershell
# Create and switch to new branch
git checkout -b feature/new-feature

# Make changes, commit them
git add .
git commit -m "Add new feature"

# Push branch to GitHub
git push -u origin feature/new-feature
```

Then create a Pull Request on GitHub to merge into main.

## 🔐 Authentication

GitHub may prompt for credentials. You have options:

### Option 1: Personal Access Token (Recommended)
1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate new token with `repo` permissions
3. Use token as password when pushing

### Option 2: GitHub CLI
Install GitHub CLI: https://cli.github.com/
```powershell
gh auth login
```

### Option 3: SSH Keys
Set up SSH keys: https://docs.github.com/en/authentication/connecting-to-github-with-ssh

## 🚨 Important Notes

1. **Never commit .env files** - They contain secrets! (Already in .gitignore)
2. **Review files before committing** - Use `git status` to see what will be committed
3. **Write meaningful commit messages** - Describe what changed and why
4. **Keep README updated** - Update documentation as you add features

## 📋 Quick Reference

```powershell
# Check status
git status

# See changes
git diff

# View commit history
git log --oneline

# Undo last commit (keep changes)
git reset --soft HEAD~1

# Discard local changes
git checkout -- filename

# Pull latest changes
git pull
```

## 🎯 Next Steps After Push

1. Add GitHub repository URL to your README
2. Set up GitHub Actions (CI/CD) - optional
3. Enable GitHub Issues for bug tracking
4. Add topics/tags to your repository
5. Create a LICENSE file if public
6. Add repository description and website URL

---

**Your project is ready to share with the world! 🚀**
