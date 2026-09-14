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

cd /opt/portfolio

# Fail fast, before anything gets torn down: this script overwrites
# nginx/conf.d/default.conf with a temporary HTTP-only config in step 1 and
# only replaces it with a working HTTPS config in step 4. If the production
# template is missing, checking that now avoids leaving the site stuck on
# "SSL setup in progress..." with no way back to serving HTTPS.
if [ ! -f nginx/templates/default.prod.conf.template ]; then
    error "nginx/templates/default.prod.conf.template not found"
    error "This template is tracked in git; make sure your deployment checkout is up to date"
    exit 1
fi

info "Setting up SSL certificate for: $DOMAIN"
info "Email: $SSL_EMAIL"

# ============================================
# 1. Create temporary nginx config (HTTP only)
# ============================================
info "Creating temporary nginx configuration..."

cat > nginx/conf.d/default.conf << EOF
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
# 2. Restart nginx with temporary config
# ============================================
info "Restarting nginx with temporary config..."
docker compose restart nginx
sleep 5
success "Nginx restarted"

# ============================================
# 3. Obtain SSL certificate
# ============================================
info "Obtaining SSL certificate from Let's Encrypt..."
docker run --rm \
    --name certbot-manual \
    -v "$(pwd)/certbot/www:/var/www/certbot:rw" \
    -v "$(pwd)/certbot/conf:/etc/letsencrypt:rw" \
    certbot/certbot:latest \
    certonly \
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

# The production config is tracked in git as a template (see
# nginx/templates/default.prod.conf.template) instead of being generated
# inline here, so the config actually served in production can be reviewed
# and diffed like any other file. This just substitutes the domain into it.
# (Existence of the template was already checked at the top of this script.)
sed "s/DOMAIN_PLACEHOLDER/$DOMAIN/g" nginx/templates/default.prod.conf.template > nginx/conf.d/default.conf

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
success "SSL setup completed!"
echo ""
info "Your site is now accessible at:"
echo "  https://$DOMAIN"
echo "  https://www.$DOMAIN"
echo ""
info "SSL certificate will auto-renew via certbot container"
echo ""
success "Deployment complete!"
