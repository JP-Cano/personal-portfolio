# 🚀 Deployment Guide - DigitalOcean with Docker & GitHub Actions

This guide will help you deploy your portfolio to DigitalOcean with automatic CI/CD using GitHub Actions.

## 📋 Prerequisites

- DigitalOcean account
- GitHub account
- Domain name (optional but recommended)
- SSH key pair

## 🎯 Overview

This setup includes:
- ✅ Docker & Docker Compose for containerization
- ✅ Nginx as reverse proxy
- ✅ Let's Encrypt SSL certificates (auto-renewal)
- ✅ GitHub Actions for CI/CD (auto-deploy on push to main)
- ✅ Zero-downtime deployments
- ✅ Health checks for all services

---

## 🌊 Step 1: Create DigitalOcean Droplet

1. **Log in to DigitalOcean** and create a new Droplet:
   - **Image**: Ubuntu 24.04 LTS
   - **Plan**: Basic ($6/month - 1GB RAM, 1 vCPU)
   - **Datacenter**: Choose closest to you
   - **Authentication**: SSH Key (add your public key)
   - **Hostname**: `portfolio-server` (or your choice)

2. **Note your Droplet IP address** - you'll need it later

---

## 🔧 Step 2: Set Up Your Droplet

### Connect to your Droplet:
```bash
ssh root@YOUR_DROPLET_IP
```

### Run the setup script:
```bash
# Download the setup script
curl -o setup-droplet.sh https://raw.githubusercontent.com/YOUR_USERNAME/personal-portfolio/main/scripts/setup-droplet.sh

# Make it executable
chmod +x setup-droplet.sh

# Run it
sudo bash setup-droplet.sh
```

This script will:
- Update system packages
- Install Docker & Docker Compose
- Create deployment user (`deployer`)
- Set up firewall (UFW)
- Create application directories
- Configure security updates

---

## 🔑 Step 3: Set Up SSH Key for Deployment

### On your local machine, generate a deployment key:
```bash
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/portfolio_deploy
```

### Add the public key to your Droplet:
```bash
# Copy your public key
cat ~/.ssh/portfolio_deploy.pub

# SSH into your droplet as root
ssh root@YOUR_DROPLET_IP

# Add the key to deployer user
echo "YOUR_PUBLIC_KEY_HERE" >> /home/deployer/.ssh/authorized_keys

# Test the connection
ssh -i ~/.ssh/portfolio_deploy deployer@YOUR_DROPLET_IP
```

---

## 📦 Step 4: Deploy Project Files to Droplet

### From your local machine:
```bash
# Copy required files to droplet
scp -r docker-compose.yml nginx scripts deployer@YOUR_DROPLET_IP:/opt/portfolio/

# Copy environment file
cp .env.example .env
# Edit .env with your configuration
nano .env

# Copy .env to droplet
scp .env deployer@YOUR_DROPLET_IP:/opt/portfolio/
```

### Your `.env` file should include:
```env
# Backend
PORT=8080
DEBUG=false
DB_DRIVER=sqlite
DATABASE_PATH=portfolio.db

# Frontend
FRONTEND_PORT=4321
PORTFOLIO_BACKEND_URL=http://backend:8080
ENVIRONMENT=production

# Domain (update these!)
DOMAIN=yourdomain.com
SSL_EMAIL=your-email@example.com
```

---

## 🌐 Step 5: Configure Your Domain (Optional)

If you have a domain:

1. **Add DNS A Records** pointing to your Droplet IP:
   ```
   A     @              YOUR_DROPLET_IP
   A     www            YOUR_DROPLET_IP
   ```

2. **Wait for DNS propagation** (can take up to 48 hours, usually much faster)

3. **Verify DNS**:
   ```bash
   dig yourdomain.com
   ```

---

## 🔒 Step 6: Set Up SSL Certificate

### SSH into your Droplet:
```bash
ssh deployer@YOUR_DROPLET_IP
cd /opt/portfolio
```

### Run SSL setup script:
```bash
chmod +x scripts/ssl-setup.sh
bash scripts/ssl-setup.sh
```

This will:
- Start nginx with HTTP configuration
- Obtain SSL certificates from Let's Encrypt
- Update nginx with HTTPS configuration
- Start all services with SSL

### Verify SSL:
```bash
# Check if all services are running
docker compose -f docker-compose.yml ps

# Test HTTPS
curl -I https://yourdomain.com
```

---

## 🤖 Step 7: Set Up GitHub Actions (CI/CD)

### 1. Make your GitHub Container Registry (GHCR) public or set up access

In your GitHub repository settings:
- Go to **Settings** → **Packages** → **Manage Action Permissions**
- Ensure actions have read/write access

### 2. Add Secrets to GitHub Repository

Go to your GitHub repository → **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

Add these secrets:
```
DROPLET_HOST = YOUR_DROPLET_IP
DROPLET_USERNAME = deployer
DROPLET_SSH_KEY = (paste content of ~/.ssh/portfolio_deploy private key)
```

### 3. Test the GitHub Actions workflow

```bash
# Make a small change to your code
echo "# Deployment test" >> README.md

# Commit and push
git add .
git commit -m "test: trigger deployment"
git push origin main
```

### 4. Monitor deployment

- Go to your repository → **Actions** tab
- Watch the deployment workflow run
- Check your site: `https://yourdomain.com`

---

## 📊 Step 8: Verify Deployment

### Check all services:
```bash
ssh deployer@YOUR_DROPLET_IP
cd /opt/portfolio

# View running containers
docker compose -f docker-compose.yml ps

# View logs
docker compose -f docker-compose.yml logs -f

# Check individual service logs
docker compose -f docker-compose.yml logs backend
docker compose -f docker-compose.yml logs frontend
docker compose -f docker-compose.yml logs nginx
```

### Test your endpoints:
```bash
# Health check
curl http://YOUR_DROPLET_IP/health

# API
curl https://yourdomain.com/api/v1/health

# Frontend
curl https://yourdomain.com
```

---

## 🔄 Continuous Deployment

From now on, every time you push to the `main` branch:

1. ✅ GitHub Actions builds Docker images
2. ✅ Images are pushed to GitHub Container Registry
3. ✅ SSH into your Droplet
4. ✅ Pull new images
5. ✅ Restart services with zero downtime
6. ✅ Clean up old images

**You don't need to do anything manually!** Just push your code.

---

## 🛠️ Common Commands

### On your Droplet:

```bash
# SSH into droplet
ssh deployer@YOUR_DROPLET_IP
cd /opt/portfolio

# View all containers
docker compose -f docker-compose.yml ps

# View logs (all services)
docker compose -f docker-compose.yml logs -f

# View specific service logs
docker compose -f docker-compose.yml logs -f backend
docker compose -f docker-compose.yml logs -f frontend

# Restart all services
docker compose -f docker-compose.yml restart

# Restart specific service
docker compose -f docker-compose.yml restart backend

# Stop all services
docker compose -f docker-compose.yml down

# Start all services
docker compose -f docker-compose.yml up -d

# Rebuild and restart
docker compose -f docker-compose.yml up -d --build

# Check disk usage
docker system df

# Clean up unused images/containers
docker system prune -a

# View resource usage
docker stats
```

---

## 🗄️ Database Management

### SQLite (default):
Your database is persisted in `/opt/portfolio/data/backend/portfolio.db`

### Backup database:
```bash
# SSH into droplet
ssh deployer@YOUR_DROPLET_IP

# Backup SQLite database
docker compose -f /opt/portfolio/docker-compose.yml exec backend \
  cp /root/portfolio.db /root/portfolio.db.backup

# Copy backup to local machine
scp deployer@YOUR_DROPLET_IP:/opt/portfolio/data/backend/portfolio.db.backup ./portfolio-backup-$(date +%Y%m%d).db
```

### Migrate to Turso (optional):
If you want to use [Turso](https://turso.tech) (distributed SQLite):

1. Sign up at https://turso.tech
2. Create a database
3. Get your database URL and auth token
4. Update `.env` on your Droplet:
   ```env
   DB_DRIVER=turso
   TURSO_DATABASE_URL=libsql://your-db.turso.io
   TURSO_AUTH_TOKEN=your-token-here
   ```
5. Restart backend:
   ```bash
   docker compose -f docker-compose.yml restart backend
   ```

---

## 🔐 Security Best Practices

### ✅ Already configured:
- Firewall (UFW) allows only SSH, HTTP, HTTPS
- Non-root deployment user
- Automatic security updates
- SSL/TLS with Let's Encrypt
- Security headers in nginx

### 🔒 Additional recommendations:

1. **Disable root SSH login**:
   ```bash
   sudo nano /etc/ssh/sshd_config
   # Set: PermitRootLogin no
   sudo systemctl restart ssh
   ```

2. **Set up fail2ban** (prevents brute-force attacks):
   ```bash
   sudo apt-get install fail2ban
   sudo systemctl enable fail2ban
   sudo systemctl start fail2ban
   ```

3. **Enable DigitalOcean Monitoring**:
   - Go to your Droplet → Monitoring
   - Enable monitoring and set up alerts

4. **Regular updates**:
   ```bash
   # Update system packages weekly
   sudo apt-get update && sudo apt-get upgrade -y
   
   # Update Docker images
   cd /opt/portfolio
   docker compose -f docker-compose.yml pull
   docker compose -f docker-compose.yml up -d
   ```

---

## 📈 Monitoring & Logs

### View application logs:
```bash
# All services
docker compose -f docker-compose.yml logs -f

# Last 100 lines
docker compose -f docker-compose.yml logs --tail=100

# Specific service
docker compose -f docker-compose.yml logs -f backend
```

### Monitor resource usage:
```bash
# Real-time stats
docker stats

# System resources
htop

# Disk usage
df -h
docker system df
```

### Set up log rotation:
Docker automatically rotates logs, but you can configure it:

```bash
# Edit Docker daemon config
sudo nano /etc/docker/daemon.json
```

Add:
```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

```bash
# Restart Docker
sudo systemctl restart docker
```

---

## 🐛 Troubleshooting

### Service won't start:
```bash
# Check logs
docker compose -f docker-compose.yml logs backend

# Check if port is in use
sudo netstat -tulpn | grep :8080

# Restart service
docker compose -f docker-compose.yml restart backend
```

### SSL certificate issues:
```bash
# Check certificate status
docker compose -f docker-compose.yml run --rm certbot certificates

# Renew certificate manually
docker compose -f docker-compose.yml run --rm certbot renew

# Restart nginx
docker compose -f docker-compose.yml restart nginx
```

### Out of disk space:
```bash
# Check disk usage
df -h
docker system df

# Clean up
docker system prune -a --volumes
docker image prune -a
```

### GitHub Actions deployment fails:
- Check Actions logs in GitHub
- Verify secrets are set correctly
- Ensure SSH key has correct permissions
- Check if services are running on Droplet

---

## 💰 Cost Optimization

Your current setup costs **$6/month** with DigitalOcean's basic droplet.

### Ways to reduce costs:

1. **Use Turso for database** (free tier: 9GB storage, 500M row reads/month)
2. **Optimize Docker images** (already using multi-stage builds)
3. **Use CDN for static assets** (Cloudflare free tier)
4. **Monitor resource usage** and downsize if possible

---

## 🎓 What You've Learned

By completing this deployment, you've gained experience with:

✅ **Docker & Docker Compose** - containerization and orchestration  
✅ **Nginx** - reverse proxy and web server configuration  
✅ **SSL/TLS** - Let's Encrypt certificate management  
✅ **GitHub Actions** - CI/CD pipeline automation  
✅ **Linux server administration** - user management, firewall, SSH  
✅ **DevOps practices** - zero-downtime deployments, health checks  
✅ **Infrastructure as Code** - reproducible deployments  

---

## 🚀 Next Steps

1. **Custom domain**: If you haven't already, configure your domain
2. **Monitoring**: Set up monitoring/alerting (Uptime Robot, Better Stack)
3. **Analytics**: Add analytics to your portfolio (Plausible, Simple Analytics)
4. **Performance**: Enable caching in nginx for static assets
5. **Backups**: Set up automated database backups
6. **Staging environment**: Create a separate branch for testing

---

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [DigitalOcean Tutorials](https://www.digitalocean.com/community/tutorials)

---

## 🆘 Need Help?

If you run into issues:

1. Check the troubleshooting section above
2. Review Docker logs: `docker compose -f docker-compose.yml logs`
3. Check GitHub Actions logs in your repository
4. Verify all environment variables in `.env`
5. Ensure DNS is properly configured for your domain

---

## 🎉 Congratulations!

You've successfully deployed your portfolio with a professional CI/CD pipeline! Every time you push to `main`, your changes will automatically deploy to production.

**Happy coding!** 🚀
