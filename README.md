# 🚀 Personal Portfolio

> A modern, full-stack web application showcasing professional experience, projects, and certifications with a beautiful 3D interactive interface.

[![Deploy](https://img.shields.io/badge/DigitalOcean-0080FF?style=flat&logo=digitalocean&logoColor=white)](https://github.com/JuanPabloCano/personal-portfolio/actions)
[![Go Version](https://img.shields.io/badge/Go-1.25.1-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Astro](https://img.shields.io/badge/Astro-5.16.5-FF5D01?style=flat&logo=astro)](https://astro.build/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

---

## ✨ Features

- 🎨 **Beautiful 3D Interactive Interface** - Three.js powered background with smooth animations
- ⚡ **Lightning Fast** - Built with Astro for optimal performance and SEO
- 🔒 **Secure Backend API** - RESTful API built with Go and Gin framework
- 📱 **Fully Responsive** - Seamless experience across all devices
- 🌓 **Dark/Light Mode** - Theme switching with system preference detection
- 🔐 **Session-Based Authentication** - Secure admin panel with cookie-based sessions
- 📊 **Swagger Documentation** - Interactive API documentation
- 🐳 **Docker Ready** - Containerized application with Docker Compose
- 🚀 **CI/CD Pipeline** - Automated deployments with GitHub Actions
- 🔄 **Zero-Downtime Deployments** - Health checks and rolling updates
- 🔒 **SSL/TLS** - Automated certificate management with Let's Encrypt

---

## 🏗️ Architecture

This is a modern full-stack application with a clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                         Nginx (Reverse Proxy)                │
│                    SSL/TLS + Load Balancing                  │
└───────────────┬─────────────────────────┬───────────────────┘
                │                         │
        ┌───────▼────────┐        ┌───────▼────────┐
        │   Frontend     │        │    Backend     │
        │   (Astro)      │───────▶│     (Go)       │
        │   Port: 4321   │        │   Port: 8080   │
        └────────────────┘        └────────┬───────┘
                                           │
                                   ┌───────▼────────┐
                                   │   Database     │
                                   │ SQLite/Turso   │
                                   └────────────────┘
```

---

## 🛠️ Tech Stack

### Frontend
- **[Astro](https://astro.build/)** - Modern web framework for content-focused sites
- **[Three.js](https://threejs.org/)** - 3D graphics library for interactive backgrounds
- **[GLightbox](https://github.com/biati-digital/glightbox)** - Responsive lightbox gallery
- **[TypeScript](https://www.typescriptlang.org/)** - Type-safe JavaScript
- **[Zod](https://zod.dev/)** - Schema validation

### Backend
- **[Go 1.25](https://go.dev/)** - High-performance compiled language
- **[Gin](https://gin-gonic.com/)** - Fast HTTP web framework
- **[GORM](https://gorm.io/)** - ORM library for Go
- **[Swagger](https://swagger.io/)** - API documentation
- **[SQLite](https://www.sqlite.org/)** / **[Turso](https://turso.tech/)** - Database options
- **[Goose](https://github.com/pressly/goose)** - Database migrations
- **[UUID](https://github.com/google/uuid)** - Session ID generation

### DevOps & Infrastructure
- **[Docker](https://www.docker.com/)** - Containerization
- **[Docker Compose](https://docs.docker.com/compose/)** - Multi-container orchestration
- **[Nginx](https://nginx.org/)** - Reverse proxy and load balancer
- **[GitHub Actions](https://github.com/features/actions)** - CI/CD automation
- **[Let's Encrypt](https://letsencrypt.org/)** - Free SSL/TLS certificates
- **[DigitalOcean](https://www.digitalocean.com/)** - Cloud hosting platform

---

## 📁 Project Structure

```
personal-portfolio/
├── backend/                    # Go backend application
│   ├── cmd/api/               # Application entry point
│   ├── internal/              # Private application code
│   │   ├── handlers/          # HTTP request handlers
│   │   ├── middleware/        # Custom middleware
│   │   ├── models/            # Data models
│   │   ├── repository/        # Data access layer
│   │   ├── routes/            # Route definitions
│   │   └── services/          # Business logic
│   ├── pkg/                   # Public packages
│   │   ├── database/          # Database connection
│   │   ├── logger/            # Logging utilities
│   │   └── utils/             # Helper functions
│   ├── migrations/            # Database migrations
│   ├── docs/                  # Swagger documentation
│   └── Dockerfile             # Backend container image
│
├── frontend/                   # Astro frontend application
│   ├── src/
│   │   ├── components/        # Reusable components
│   │   ├── features/          # Feature sections
│   │   ├── layouts/           # Page layouts
│   │   ├── pages/             # Route pages
│   │   ├── icons/             # SVG icon components
│   │   ├── api/               # API client
│   │   ├── types/             # TypeScript types
│   │   └── utils/             # Utility functions
│   ├── public/                # Static assets
│   └── Dockerfile             # Frontend container image
│
├── nginx/                      # Nginx configuration
│   ├── nginx.conf             # Main config
│   └── conf.d/                # Site-specific configs
│
├── scripts/                    # Deployment scripts
│   ├── setup-droplet.sh       # Server setup
│   ├── ssl-setup.sh           # SSL configuration
│   └── local-test.sh          # Local testing
│
├── docs/                       # Documentation
│   ├── DEPLOYMENT.md          # Deployment guide
│   ├── DOCKER-SETUP.md        # Docker setup guide
│   └── CORS-GUIDE.md          # CORS configuration
│
├── .github/workflows/         # CI/CD pipelines
│   └── deploy.yml             # Auto-deployment workflow
│
├── docker-compose.yml         # Multi-container setup
└── README.md                  # This file
```

---

## 🔌 API Endpoints

### Public Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/projects` | Get all projects |
| GET | `/api/v1/projects/:id` | Get project by ID |
| GET | `/api/v1/experiences` | Get all experiences |
| GET | `/api/v1/experiences/:id` | Get experience by ID |
| GET | `/api/v1/certifications` | Get all certifications |
| GET | `/api/v1/certifications/:id` | Get certification by ID |

### Protected Endpoints (Admin)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Admin login |
| POST | `/api/v1/auth/logout` | Admin logout |
| POST | `/api/v1/projects` | Create project |
| PUT | `/api/v1/projects/:id` | Update project |
| DELETE | `/api/v1/projects/:id` | Delete project |
| POST | `/api/v1/experiences` | Create experience |
| PUT | `/api/v1/experiences/:id` | Update experience |
| DELETE | `/api/v1/experiences/:id` | Delete experience |
| POST | `/api/v1/certifications` | Upload certification |
| PUT | `/api/v1/certifications/:id` | Update certification |
| DELETE | `/api/v1/certifications/:id` | Delete certification |

📚 **Full API Documentation:** Available at `/api/v1/swagger/index.html`

---

## 🌐 Deployment

This project includes a complete CI/CD pipeline for automated deployments to DigitalOcean.

### Automated Deployment (GitHub Actions)

Every push to `main` branch automatically:
1. ✅ Builds Docker images for backend and frontend
2. ✅ Pushes images to GitHub Container Registry
3. ✅ Deploys to DigitalOcean droplet via SSH
4. ✅ Performs zero-downtime rolling updates
5. ✅ Cleans up old images

### Manual Deployment

For detailed deployment instructions, see [📖 DEPLOYMENT.md](docs/DEPLOYMENT.md)

**Quick deployment steps:**
```bash
# 1. Set up your server
bash scripts/setup-droplet.sh

# 2. Configure SSL
bash scripts/ssl-setup.sh

# 3. Deploy application
docker compose pull
docker compose up -d
```

---

## 🔧 Configuration

### Environment Variables

#### Backend (.env)
```env
# Server
SERVER_PORT=8080
DEBUG=false

# Database (choose one)
DB_DRIVER=sqlite                        # or "turso"
DATABASE_PATH=portfolio.db              # for SQLite
TURSO_DATABASE_URL=libsql://...         # for Turso
TURSO_AUTH_TOKEN=your-token             # for Turso

# Security
SESSION_SECRET=your-session-secret-key
ALLOWED_ORIGINS=http://localhost:4321,https://yourdomain.com

# Admin Credentials
ADMIN_EMAIL=your-email@example.com
ADMIN_PASSWORD=your-secure-password
```

#### Frontend (.env)
```env
PORTFOLIO_BACKEND_URL=http://localhost:8080/api/v1
ENVIRONMENT=development
```

---

## 📊 Database Migrations

Migrations are managed with [Goose](https://github.com/pressly/goose).

```bash
cd backend

# Create new migration
goose -dir migrations create migration_name sql

# Run migrations
make migrate-up

# Rollback
make migrate-down

# Check status
goose -dir migrations sqlite3 portfolio.db status
```

---

## 🎨 Key Features Showcase

### 🌟 3D Interactive Background
Powered by Three.js, featuring animated particle systems that respond to mouse movement, creating an engaging visual experience.

### 🔐 Secure Admin Panel
Session-based authentication with secure HTTP-only cookies, allowing safe content management through a simple email/password login.

### ⚡ Performance Optimized
- Server-side rendering with Astro
- Lazy loading for images and components
- Optimized Docker images with multi-stage builds
- Nginx caching and compression

### 📱 Mobile-First Design
Responsive design that looks great on all devices, from phones to desktop monitors.

---

## 👤 Author

**Juan Pablo Cano**

- GitHub: [@JuanPabloCano](https://github.com/JuanPabloCano)
- LinkedIn: [Juan Pablo Cano](https://linkedin.com/in/your-profile)

---

## 🙏 Acknowledgments

- [Astro](https://astro.build/) - For the amazing web framework
- [Gin](https://gin-gonic.com/) - For the fast and elegant Go framework
- [Three.js](https://threejs.org/) - For making 3D graphics accessible
- [DigitalOcean](https://www.digitalocean.com/) - For reliable hosting

---

## 📚 Additional Documentation

- [📦 Docker Setup Guide](docs/DOCKER-SETUP.md)
- [🚀 Deployment Guide](docs/DEPLOYMENT.md)
- [🔧 CORS Configuration](docs/CORS-GUIDE.md)

---
