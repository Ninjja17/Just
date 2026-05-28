# Quick Start Guide

## Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+

## Setup Instructions

### 1. Backend Setup

```powershell
# Navigate to backend directory
cd backend

# Install Go dependencies
go mod download
go mod tidy

# Create .env file from example
copy .env.example .env

# Update .env with your credentials:
# - Database connection details
# - Google OAuth credentials (get from Google Cloud Console)
# - SendGrid API key (get from SendGrid)
# - JWT secret (generate a random string)

# Run the backend (make sure PostgreSQL and Redis are running)
go run cmd/server/main.go
```

### 2. Frontend Setup

```powershell
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Create .env file from example
copy .env.example .env

# Start development server
npm run dev
```

### 3. Database Setup

Option 1: Using Docker
```powershell
docker-compose up -d postgres redis
```

Option 2: Local Installation
- Install PostgreSQL 15+
- Create database: `tenthousand_hr`
- Run migration: `backend/migrations/001_initial_schema.sql`

### 4. Full Docker Setup (All Services)

```powershell
# From project root
docker-compose up -d
```

## Access the Application

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- API Health Check: http://localhost:8080/health

## Next Steps

1. Register a new account
2. Verify your email with OTP
3. Create your first skill
4. Start tracking sessions
5. Set goals and track progress!

## Development

### Backend Development
```powershell
cd backend
go run cmd/server/main.go
```

### Frontend Development
```powershell
cd frontend
npm run dev
```

### Building for Production
```powershell
# Backend
cd backend
go build -o bin/server cmd/server/main.go

# Frontend
cd frontend
npm run build
```

## Troubleshooting

**Database connection failed:**
- Check PostgreSQL is running
- Verify credentials in .env file
- Ensure database exists

**Redis connection failed:**
- Check Redis is running
- Verify Redis host and port in .env

**OTP emails not sending:**
- Verify SendGrid API key is valid
- Check FROM_EMAIL is configured

**Google OAuth not working:**
- Verify Google Client ID and Secret
- Check redirect URL is configured in Google Cloud Console
