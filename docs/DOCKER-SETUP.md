# 🐳 Docker Setup Guide

## 📋 Overview

This project uses **Docker Compose** with different configurations for local development and production.

---

## 🏠 **Local Development (No nginx/certbot needed)**

### **Files:**
- `docker-compose.yml` - Base configuration
- `docker-compose.override.yml` - Local development overrides (exposes ports)

### **What happens locally:**
- ✅ Backend exposed on `localhost:8080`
- ✅ Frontend exposed on `localhost:4321`
- ❌ Nginx disabled (not needed)
- ❌ Certbot disabled (not needed)

### **Commands:**

```bash
# Start backend + frontend (nginx/certbot won't start)
docker compose up -d backend frontend

# Or start everything (but nginx/certbot are disabled by profile)
docker compose up -d

# View logs
docker compose logs -f

# Stop
docker compose down
```

### **Access your app:**
- Frontend: http://localhost:4321
- Backend API: http://localhost:8080/api/v1/
- Swagger: http://localhost:8080/swagger/index.html
- Health: http://localhost:8080/health

---

## 🚀 **Production Deployment (DigitalOcean)**

### **Files:**
- `docker-compose.yml` - Base configuration (NO ports exposed)
- Services accessed ONLY through nginx

### **What happens in production:**
- ✅ Nginx exposes ports 80/443
- ✅ SSL/HTTPS with Let's Encrypt
- ✅ All traffic goes through nginx reverse proxy
- ❌ Backend/frontend ports NOT exposed (security)

### **Architecture:**

```
Internet → nginx (80/443) → backend (internal 8080)
                          → frontend (internal 4321)
```

### **Commands on server:**

```bash
# On DigitalOcean droplet
cd /opt/portfolio

# Start all services (nginx + certbot included)
docker compose up -d

# View logs
docker compose logs -f

# Restart
docker compose restart
```

### **Access:**
- Everything through your domain: `https://yourdomain.com`
- No direct access to backend/frontend ports (secure!)

---

## 🔒 **Security: Why No Port Exposure in Production?**

### ❌ **Bad (Insecure):**
```yaml
backend:
  ports:
    - "8080:8080"  # ❌ Exposed to internet!
```

**Problems:**
- Anyone can access backend directly
- Bypasses nginx security/SSL
- No rate limiting
- Larger attack surface

### ✅ **Good (Secure):**
```yaml
backend:
  # No ports - only nginx can access
  networks:
    - portfolio-network

nginx:
  ports:
    - "80:80"    # Only nginx exposed
    - "443:443"  # SSL termination
```

**Benefits:**
- ✅ Single entry point (nginx)
- ✅ SSL/HTTPS encryption
- ✅ Rate limiting, security headers
- ✅ Smaller attack surface

---

## 📁 **File Structure**

```
.
├── docker-compose.yml              # Base config (production-ready)
├── docker-compose.override.yml     # Local dev overrides
├── docker-compose.local.yml        # Backup local config
│
├── backend/
│   └── Dockerfile                  # Backend image
│
├── frontend/
│   └── Dockerfile                  # Frontend image
│
└── nginx/
    ├── nginx.conf                  # Main nginx config
    └── conf.d/
        ├── default.conf.production # Production config (with SSL)
        └── local.conf              # Local config (HTTP only)
```

---

## 🔧 **How It Works**

### **Docker Compose Override**

Docker Compose automatically merges:
1. `docker-compose.yml` (base)
2. `docker-compose.override.yml` (if exists)

**Local development:**
```bash
docker compose up -d
# Loads: docker-compose.yml + docker-compose.override.yml
# Result: Ports exposed, nginx/certbot disabled
```

**Production:**
```bash
# Don't copy docker-compose.override.yml to server!
docker compose up -d
# Loads: docker-compose.yml only
# Result: No ports exposed, nginx/certbot enabled
```

---

## ⚙️ **Environment Variables**

### **Port Configuration:**

`.env` file:
```env
PORT=8080                # Backend port
FRONTEND_PORT=4321       # Frontend port
```

**Local:** Ports exposed via `docker-compose.override.yml`  
**Production:** Ports internal only

---

## 🧪 **Testing**

### **Local Test:**
```bash
# Start services
docker compose up -d backend frontend

# Test backend
curl http://localhost:8080/health

# Test frontend
curl http://localhost:4321

# Test API
curl http://localhost:8080/api/v1/health
```

### **Production Test (on server):**
```bash
# Through nginx
curl http://localhost/health
curl https://yourdomain.com

# Check services are NOT exposed
curl http://localhost:8080  # Should fail (good!)
curl http://localhost:4321  # Should fail (good!)
```

---

## 🚨 **Common Issues**

### **"Port already in use"**
```bash
# Find what's using the port
lsof -i :8080

# Stop conflicting service
docker compose down
```

### **"Can't access frontend/backend"**
```bash
# Check if services are running
docker compose ps

# Check if ports are exposed (local only)
docker compose ps | grep "0.0.0.0"

# View logs
docker compose logs backend
docker compose logs frontend
```

### **"nginx keeps restarting"**
```bash
# Check nginx logs
docker compose logs nginx

# Usually SSL certificate issue - run ssl-setup.sh
bash scripts/ssl-setup.sh
```

---

## 📊 **Port Summary**

| Service | Local Port | Production Port | Production Access |
|---------|-----------|-----------------|-------------------|
| Backend | `8080` | None (internal) | Through nginx |
| Frontend | `4321` | None (internal) | Through nginx |
| Nginx | Disabled | `80`, `443` | Public |
| Certbot | Disabled | Internal | nginx only |

---

## ✅ **Best Practices**

### **Local Development:**
1. ✅ Use `docker-compose.override.yml` for port exposure
2. ✅ Keep `.env` file (not in git)
3. ✅ Test without nginx/certbot
4. ✅ Commit `docker-compose.override.yml` (for team)

### **Production:**
1. ✅ Never expose backend/frontend ports
2. ✅ Use nginx as reverse proxy
3. ✅ Enable SSL with Let's Encrypt
4. ✅ Don't copy `docker-compose.override.yml` to server

---

## 🎯 **Quick Commands**

```bash
# Local: Start and test
docker compose up -d backend frontend
open http://localhost:4321
open http://localhost:8080/swagger/index.html

# Local: Stop
docker compose down

# Production (on server): Start all
docker compose up -d

# Production: View logs
docker compose logs -f

# Production: Restart
docker compose restart

# Check running containers
docker compose ps

# View resource usage
docker stats
```

---

## 🔍 **Debugging**

```bash
# Check if containers are healthy
docker compose ps

# View all logs
docker compose logs

# View specific service
docker compose logs -f backend

# Enter container shell
docker compose exec backend sh
docker compose exec frontend sh

# Check networks
docker network ls

# Inspect network
docker network inspect personal-portfolio_portfolio-network
```

---

## 📚 **Related Documentation**

- **DEPLOYMENT.md** - Full deployment guide
- **QUICK-START.md** - Quick reference
- **SECURITY.md** - Security audit
- **.env.example** - Environment variables

---

**Questions?** Check the troubleshooting section or review the deployment guide!

**Happy coding!** 🚀
