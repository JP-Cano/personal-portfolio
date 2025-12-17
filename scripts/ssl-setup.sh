#!/bin/bash

# ============================================
# SSL Certificate Setup with Let's Encrypt
# ============================================
# This script sets up SSL certificates using
# Let's Encrypt and Certbot for your domain
# ============================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if .env file exists
if [ ! -f /opt/portfolio/.env ]; then
    error ".env file not found in /opt/portfolio/"
    error "Please create .env file with DOMAIN and SSL_EMAIL variables"
    exit 1
fi

# Safely load environment variables (only DOMAIN and SSL_EMAIL)
# This prevents code injection by only extracting specific variables
DOMAIN=$(grep -E "^DOMAIN=" /opt/portfolio/.env | cut -d '=' -f2- | tr -d '"' | tr -d "'")
SSL_EMAIL=$(grep -E "^SSL_EMAIL=" /opt/portfolio/.env | cut -d '=' -f2- | tr -d '"' | tr -d "'")

# Check required variables
if [ -z "$DOMAIN" ] || [ -z "$SSL_EMAIL" ]; then
    error "DOMAIN and SSL_EMAIL must be set in .env file"
    exit 1
fi

info "Setting up SSL certificate for: $DOMAIN"
info "Email: $SSL_EMAIL"

# ============================================
# 1. Create temporary nginx config (HTTP only)
# ============================================
info "Creating temporary nginx configuration..."

cat > /opt/portfolio/nginx/conf.d/default.conf << EOF
server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 200 "Server is running. SSL setup in progress...";
        add_header Content-Type text/plain;
    }
}
EOF

success "Temporary nginx config created"

# ============================================
# 2. Start nginx (without certbot initially)
# ============================================
info "Starting nginx..."
cd /opt/portfolio
docker compose up -d nginx
sleep 5
success "Nginx started"

# ============================================
# 3. Obtain SSL certificate
# ============================================
info "Obtaining SSL certificate from Let's Encrypt..."
docker compose run --rm certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email "$SSL_EMAIL" \
    --agree-tos \
    --no-eff-email \
    -d "$DOMAIN" \
    -d "www.$DOMAIN"

success "SSL certificate obtained!"

# ============================================
# 4. Create production nginx config (HTTPS)
# ============================================
info "Creating production nginx configuration..."

cat > /opt/portfolio/nginx/conf.d/default.conf << 'EOF'
# HTTP server - redirects to HTTPS
server {
    listen 80;
    server_name DOMAIN_PLACEHOLDER www.DOMAIN_PLACEHOLDER;

    # Let's Encrypt challenge
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    # Health check endpoint (no redirect)
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }

    # Redirect all other HTTP traffic to HTTPS
    location / {
        return 301 https://$host$request_uri;
    }
}

# HTTPS server
server {
    listen 443 ssl http2;
    server_name DOMAIN_PLACEHOLDER www.DOMAIN_PLACEHOLDER;

    # SSL certificate files (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/DOMAIN_PLACEHOLDER/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/DOMAIN_PLACEHOLDER/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # API Backend proxy
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 90;
    }

    # Swagger documentation
    location /swagger/ {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Static files for certifications
    location /certifications/ {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Frontend (Astro) - proxy everything else
    location / {
        proxy_pass http://frontend:4321;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
EOF

# Replace domain placeholder
sed -i "s/DOMAIN_PLACEHOLDER/$DOMAIN/g" /opt/portfolio/nginx/conf.d/default.conf

success "Production nginx config created"

# ============================================
# 5. Restart nginx to apply SSL configuration
# ============================================
info "Restarting nginx with SSL configuration..."
docker compose restart nginx

success "Nginx restarted with SSL!"

# ============================================
# 6. Display summary
# ============================================
echo ""
success "🎉 SSL setup completed!"
echo ""
info "Your site is now accessible at:"
echo "  https://$DOMAIN"
echo "  https://www.$DOMAIN"
echo ""
info "SSL certificate will auto-renew via certbot container"
echo ""
success "Deployment complete! 🚀"
