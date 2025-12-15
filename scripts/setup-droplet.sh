#!/bin/bash

# ============================================
# DigitalOcean Droplet Setup Script
# ============================================
# This script sets up your DigitalOcean droplet
# for deploying your portfolio application
# ============================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
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

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    error "Please run this script as root (use: sudo bash setup-droplet.sh)"
    exit 1
fi

info "Starting DigitalOcean Droplet Setup..."

# ============================================
# 1. Update system packages
# ============================================
info "Updating system packages..."
apt-get update
apt-get upgrade -y
success "System packages updated"

# ============================================
# 2. Install Docker
# ============================================
info "Installing Docker..."
if ! command -v docker &> /dev/null; then
    # Install prerequisites
    apt-get install -y ca-certificates curl
    
    # Add Docker's official GPG key
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc
    
    # Set up the repository
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Install Docker Engine
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    # Start and enable Docker
    systemctl start docker
    systemctl enable docker
    
    # Verify Docker installation
    docker --version
    
    success "Docker installed successfully"
else
    success "Docker is already installed"
    docker --version
fi

# ============================================
# 3. Install Docker Compose
# ============================================
info "Checking Docker Compose..."
if ! docker compose version &> /dev/null; then
    error "Docker Compose plugin not found. Installing..."
    apt-get install -y docker-compose-plugin
fi
success "Docker Compose is ready"

# ============================================
# 4. Create application directory
# ============================================
info "Creating application directory..."
mkdir -p /opt/portfolio/{data/backend,data/certifications,nginx/conf.d,certbot/www,certbot/conf}
success "Application directory created"

# ============================================
# 5. Set up firewall (UFW)
# ============================================
info "Configuring firewall..."
if command -v ufw &> /dev/null; then
    ufw --force enable
    ufw default deny incoming
    ufw default allow outgoing
    ufw allow ssh
    ufw allow 80/tcp
    ufw allow 443/tcp
    ufw --force reload
    success "Firewall configured"
else
    warning "UFW not found, skipping firewall setup"
fi

# ============================================
# 6. Create non-root user for deployment
# ============================================
info "Creating deployment user..."
if ! id "deployer" &>/dev/null; then
    useradd -m -s /bin/bash deployer
    usermod -aG docker deployer
    
    # Set up SSH key for deployer (you'll need to add your key later)
    mkdir -p /home/deployer/.ssh
    touch /home/deployer/.ssh/authorized_keys
    chmod 700 /home/deployer/.ssh
    chmod 600 /home/deployer/.ssh/authorized_keys
    chown -R deployer:deployer /home/deployer/.ssh
    
    success "Deployment user 'deployer' created"
else
    success "Deployment user 'deployer' already exists"
fi

# ============================================
# 7. Set permissions
# ============================================
info "Setting permissions..."
chown -R deployer:deployer /opt/portfolio
success "Permissions set"

# ============================================
# 8. Install useful utilities
# ============================================
info "Installing useful utilities..."
apt-get install -y git wget curl vim htop
success "Utilities installed"

# ============================================
# 9. Set up automatic security updates
# ============================================
info "Setting up automatic security updates..."
apt-get install -y unattended-upgrades
dpkg-reconfigure -plow unattended-upgrades
success "Automatic security updates enabled"

# ============================================
# 10. Display summary
# ============================================
echo ""
success "DigitalOcean Droplet setup completed!"
echo ""
info "Next steps:"
echo "  1. Add your SSH key to /home/deployer/.ssh/authorized_keys"
echo "  2. Copy your project files to /opt/portfolio/"
echo "  3. Copy .env file with your configuration"
echo "  4. Update nginx/conf.d/default.conf with your domain"
echo "  5. Set up SSL certificate (see ssl-setup.sh)"
echo "  6. Configure GitHub Actions secrets:"
echo "     - DROPLET_HOST: your droplet IP"
echo "     - DROPLET_USERNAME: deployer"
echo "     - DROPLET_SSH_KEY: your private SSH key"
echo ""
info "To copy files from local machine, use:"
echo "  scp -r docker-compose.yml nginx .env deployer@YOUR_DROPLET_IP:/opt/portfolio/"
echo ""
