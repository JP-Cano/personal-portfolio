# 🌐 CORS Configuration Guide

## Understanding CORS in Different Environments

### 🔍 **The Problem**

CORS (Cross-Origin Resource Sharing) is enforced by browsers when:
- Frontend origin: `http://localhost:3000`
- Backend origin: `http://localhost:8080`

Different **ports** = Different **origins** = CORS check required!

---

## 📋 **Configuration by Environment**

### **1️⃣ Local Development (No Docker)**

**How it works:**
```
Browser → Frontend (localhost:4321) → Backend (localhost:8080)
```

**`.env` configuration:**
```env
# Backend will allow these origins
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4321

# Frontend will call backend at
PORTFOLIO_BACKEND_URL=http://localhost:8080
```

**Why both ports?**
- `4321` - Astro dev server (`pnpm run dev`)
- `3000` - Astro preview server (`pnpm preview`)

---

### **2️⃣ Docker Local Testing**

**How it works:**
```
Browser (localhost:3000) → Frontend Container (exposed on :3000)
                         ↓
                    Backend Container (exposed on :8080)
                         ↑
                    Both accessible via localhost from browser
```

**Key Point:** Your **browser** runs on your host machine, not inside Docker!

**`.env` configuration:**
```env
# Backend allows requests from browser
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4321

# Frontend tells browser where backend is (on host machine)
PORTFOLIO_BACKEND_URL=http://localhost:8080
```

**Why `localhost:8080` and not `backend:8080`?**
- ❌ `http://backend:8080` - Only works inside Docker network
- ✅ `http://localhost:8080` - Accessible from browser on host machine

---

### **3️⃣ Production (DigitalOcean with Nginx)**

**How it works:**
```
Browser (yourdomain.com) → Nginx (:443) → Frontend Container
                                       ↓
                                  Backend Container
                                  (internal only)
```

**`.env` on server:**
```env
# Backend allows requests from your domain
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Frontend calls backend internally through Docker network
PORTFOLIO_BACKEND_URL=http://backend:8080
```

**Why `http://backend:8080` in production?**
- Nginx handles SSL termination
- Nginx proxies requests to backend internally
- No CORS issues because Nginx acts as reverse proxy

---

## 🔧 **How to Fix CORS Issues**

### **Scenario 1: CORS error in local Docker**

**Error:**
```
Access to fetch at 'http://localhost:8080/api/v1/projects' from origin 'http://localhost:3000' 
has been blocked by CORS policy
```

**Solution:**

1. **Check `.env` has correct ALLOWED_ORIGINS:**
   ```env
   ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4321
   ```

2. **Check frontend calls localhost:8080:**
   ```env
   PORTFOLIO_BACKEND_URL=http://localhost:8080
   ```

3. **Rebuild containers:**
   ```bash
   docker compose down
   docker compose build --no-cache frontend
   docker compose up -d
   ```

---

### **Scenario 2: Frontend can't reach backend in Docker**

**Error in browser console:**
```
GET http://backend:8080/api/v1/health net::ERR_NAME_NOT_RESOLVED
```

**Cause:** Frontend is trying to call `http://backend:8080` from browser, but `backend` is only a Docker internal hostname.

**Solution:** Update `.env`:
```env
PORTFOLIO_BACKEND_URL=http://localhost:8080
```

Then rebuild:
```bash
docker compose down
docker compose build --no-cache frontend
docker compose up -d
```

---

### **Scenario 3: Backend not logging allowed origins**

**Check backend logs:**
```bash
docker compose logs backend | grep -i origin
```

**Expected output:**
```
Allowed origin: http://localhost:3000
Allowed origin: http://localhost:4321
```

**If empty or wrong:**

1. Verify `.env` exists and has `ALLOWED_ORIGINS`
2. Check `docker-compose.yml` passes the variable:
   ```yaml
   environment:
     - ALLOWED_ORIGINS=${ALLOWED_ORIGINS:-http://localhost:3000}
   ```

3. Restart:
   ```bash
   docker compose restart backend
   ```

---

## 📝 **Configuration Cheat Sheet**

### **Local Development (No Docker):**
```env
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4321
PORTFOLIO_BACKEND_URL=http://localhost:8080
```

### **Docker Local:**
```env
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4321
PORTFOLIO_BACKEND_URL=http://localhost:8080
```

### **Production:**
```env
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
PORTFOLIO_BACKEND_URL=http://backend:8080
```

---

## 🧪 **Testing CORS**

### **Test 1: Check backend is receiving ALLOWED_ORIGINS**
```bash
docker compose exec backend printenv ALLOWED_ORIGINS
```

**Expected:**
```
http://localhost:3000,http://localhost:4321
```

---

### **Test 2: Check CORS headers**
```bash
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: GET" \
     -X OPTIONS \
     http://localhost:8080/api/v1/health -v
```

**Look for:**
```
< Access-Control-Allow-Origin: http://localhost:3000
< Access-Control-Allow-Methods: GET,POST,PUT,PATCH,DELETE,HEAD
```

---

### **Test 3: Check frontend can reach backend**

Open browser console at `http://localhost:3000` and run:
```javascript
fetch('http://localhost:8080/api/v1/health')
  .then(r => r.json())
  .then(console.log)
  .catch(console.error)
```

**Expected:**
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

---

## 🐛 **Common Mistakes**

### ❌ **Wrong:**
```env
# Production URLs in local development
ALLOWED_ORIGINS=https://yourdomain.com

# Docker hostname from browser
PORTFOLIO_BACKEND_URL=http://backend:8080  # (in Docker local)

# Missing protocol
ALLOWED_ORIGINS=localhost:3000

# Single origin as comma-separated
ALLOWED_ORIGINS="http://localhost:3000, http://localhost:4321"  # (spaces!)
```

### ✅ **Correct:**
```env
# Local with both ports
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4321

# Localhost from browser
PORTFOLIO_BACKEND_URL=http://localhost:8080

# With protocol
ALLOWED_ORIGINS=http://localhost:3000

# No spaces in comma-separated list
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4321
```

---

## 🚀 **Quick Fix Commands**

```bash
# Rebuild everything
docker compose down
docker compose build --no-cache
docker compose up -d

# Check backend env vars
docker compose exec backend printenv | grep ALLOWED

# Check frontend env vars
docker compose exec frontend printenv | grep PORTFOLIO

# Test CORS
curl -H "Origin: http://localhost:3000" \
     -X OPTIONS \
     http://localhost:8080/api/v1/health -v

# View logs
docker compose logs -f backend
docker compose logs -f frontend
```

---

## 📚 **Summary**

| Environment | Frontend Origin | Backend URL | ALLOWED_ORIGINS |
|-------------|-----------------|-------------|-----------------|
| **No Docker** | `localhost:4321` | `localhost:8080` | `http://localhost:3000,http://localhost:4321` |
| **Docker Local** | `localhost:3000` (browser) | `localhost:8080` (exposed) | `http://localhost:3000,http://localhost:4321` |
| **Production** | `yourdomain.com` | `backend:8080` (internal) | `https://yourdomain.com,https://www.yourdomain.com` |

**Key takeaway:** In Docker locally, your browser still runs on the host machine, so use `localhost:8080` for the backend URL!
