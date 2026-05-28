# 10,000 Hour Tracker ⏱️

> A full-stack web application for tracking your journey to mastery through the 10,000-hour rule.

[![React](https://img.shields.io/badge/React-18-blue.svg)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.3-blue.svg)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.21-00ADD8.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791.svg)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**Repository:** [https://github.com/Ninjja17/Just](https://github.com/Ninjja17/Just)

Track practice sessions, set goals, visualize progress, and connect with others on their mastery journey.

## Features

- 🎯 **Multi-Skill Tracking** - Track progress across multiple skills simultaneously
- ⏱️ **Time Tracking** - Built-in timer and manual session entry
- 📊 **Progress Visualization** - Charts, heatmaps, and milestone badges
- 🎨 **Goal Setting** - Daily, weekly, and custom milestone goals
- 🔔 **Notifications** - Stay motivated with reminders and achievements
- 👥 **Social Features** - Follow friends, share progress, compete on leaderboards
- 🔐 **Secure Authentication** - Email OTP and Google OAuth login

## Tech Stack

### Backend
- **Go** (Golang) with Gin framework
- **PostgreSQL** for data storage
- **Redis** for caching and sessions
- **JWT** for authentication
- **Google OAuth2** for social login

### Frontend
- **React** with TypeScript
- **Tailwind CSS** for styling
- **Recharts** for data visualization
- **React Router** for navigation
- **Zustand** for state management

## Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Quick Start with Docker

1. Clone the repository
```bash
cd 1000hr
```

2. Set up environment variables
```bash
# Backend
cp backend/.env.example backend/.env

# Frontend
cp frontend/.env.example frontend/.env
```

3. Update environment variables with your credentials:
   - Google OAuth credentials
   - SendGrid API key for emails
   - JWT secret

4. Start all services
```bash
docker-compose up -d
```

5. Access the application
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080

### Local Development

#### Backend

```bash
cd backend

# Install dependencies
go mod download

# Run migrations
go run cmd/migrate/main.go

# Start server
go run cmd/server/main.go
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev
```

## Project Structure

```
1000hr/
├── backend/              # Go backend
│   ├── cmd/             # Application entrypoints
│   ├── internal/        # Private application code
│   ├── pkg/             # Public libraries
│   └── migrations/      # Database migrations
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # Reusable UI components
│   │   ├── pages/       # Page components
│   │   ├── services/    # API integration
│   │   └── context/     # State management
│   └── public/          # Static assets
└── docker-compose.yml   # Docker orchestration
```

## API Documentation

API endpoints are available at `http://localhost:8080/api`

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login with email/password
- `POST /api/auth/google` - Google OAuth login
- `POST /api/auth/verify-otp` - Verify OTP code

### Skills
- `GET /api/skills` - Get user's skills
- `POST /api/skills` - Create new skill
- `PUT /api/skills/:id` - Update skill
- `DELETE /api/skills/:id` - Delete skill

### Sessions
- `GET /api/sessions` - Get practice sessions
- `POST /api/sessions` - Create session
- `POST /api/sessions/start` - Start timer
- `POST /api/sessions/stop` - Stop timer

[View full API documentation](./backend/API.md)

## Database Schema

See [migrations](./backend/migrations/) for the complete database schema.

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please read [GITHUB_GUIDE.md](GITHUB_GUIDE.md) for more details on our code of conduct and development process.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 💬 Support

- 📫 For issues and questions, please [open an issue](../../issues)
- 📖 Check [SETUP.md](SETUP.md) for detailed setup instructions
- 📝 Review [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) for architecture overview

## 🙏 Acknowledgments

- Built with [React](https://reactjs.org/), [Go](https://golang.org/), and [PostgreSQL](https://www.postgresql.org/)
- Styled with [Tailwind CSS](https://tailwindcss.com/)
- Icons by [Heroicons](https://heroicons.com/)

---

**Star ⭐ this repository if you find it helpful!**
