#!/bin/bash

# EatRight Backend Deployment Script
# This script deploys the built binary to the production environment

set -e  # Exit on error

echo "🚀 Starting EatRight Backend Deployment..."

# Configuration
APP_NAME="eatright-server"
DEPLOY_DIR="/var/www/eatright"
SERVICE_NAME="eatright"
BACKUP_DIR="${DEPLOY_DIR}/backups"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Create necessary directories
log_info "Creating deployment directories..."
sudo mkdir -p ${DEPLOY_DIR}/bin
sudo mkdir -p ${BACKUP_DIR}
sudo mkdir -p ${DEPLOY_DIR}/logs

# Backup current binary if exists
if [ -f "${DEPLOY_DIR}/bin/${APP_NAME}" ]; then
    BACKUP_FILE="${BACKUP_DIR}/${APP_NAME}.$(date +%Y%m%d_%H%M%S).backup"
    log_info "Backing up current binary to ${BACKUP_FILE}..."
    sudo cp ${DEPLOY_DIR}/bin/${APP_NAME} ${BACKUP_FILE}
    
    # Keep only last 5 backups
    log_info "Cleaning old backups..."
    cd ${BACKUP_DIR}
    ls -t ${APP_NAME}.*.backup | tail -n +6 | xargs -r rm --
    cd -
fi

# Stop the service
log_info "Stopping ${SERVICE_NAME} service..."
sudo systemctl stop ${SERVICE_NAME} || log_warn "Service was not running"

# Copy new binary
log_info "Copying new binary to ${DEPLOY_DIR}/bin/..."
sudo cp bin/${APP_NAME} ${DEPLOY_DIR}/bin/${APP_NAME}
sudo chmod +x ${DEPLOY_DIR}/bin/${APP_NAME}

# Copy other necessary files
log_info "Copying configuration files..."
if [ -f "deploy/${SERVICE_NAME}.service" ]; then
    sudo cp deploy/${SERVICE_NAME}.service /etc/systemd/system/
    sudo systemctl daemon-reload
fi

if [ -f "deploy/nginx.conf" ]; then
    sudo cp deploy/nginx.conf /etc/nginx/sites-available/${SERVICE_NAME}
    sudo ln -sf /etc/nginx/sites-available/${SERVICE_NAME} /etc/nginx/sites-enabled/
fi

# Set correct permissions
log_info "Setting permissions..."
sudo chown -R eatright:eatright ${DEPLOY_DIR}

# Verify .env file exists
if [ ! -f "${DEPLOY_DIR}/.env" ]; then
    log_error ".env file not found at ${DEPLOY_DIR}/.env"
    log_error "Please create .env file before deployment!"
    exit 1
fi

# Start the service
log_info "Starting ${SERVICE_NAME} service..."
sudo systemctl start ${SERVICE_NAME}
sudo systemctl enable ${SERVICE_NAME}

# Wait a bit for service to start
sleep 3

# Check service status
if sudo systemctl is-active --quiet ${SERVICE_NAME}; then
    log_info "✅ Service is running!"
else
    log_error "❌ Service failed to start"
    log_error "Check logs with: sudo journalctl -u ${SERVICE_NAME} -n 50"
    exit 1
fi

# Reload Nginx
log_info "Reloading Nginx..."
sudo nginx -t && sudo systemctl reload nginx || log_warn "Nginx reload failed"

# Health check
log_info "Performing health check..."
sleep 2
if curl -sf http://localhost:8080/health > /dev/null; then
    log_info "✅ Health check passed!"
else
    log_warn "⚠️  Health check failed, but service is running"
fi

# Show service status
log_info "Service status:"
sudo systemctl status ${SERVICE_NAME} --no-pager -l

# Deployment summary
echo ""
echo "================================"
log_info "🎉 Deployment completed successfully!"
echo "================================"
log_info "Application: ${APP_NAME}"
log_info "Deploy directory: ${DEPLOY_DIR}"
log_info "Service: ${SERVICE_NAME}"
log_info "Health check: http://localhost:8080/health"
echo ""
log_info "View logs: sudo journalctl -u ${SERVICE_NAME} -f"
log_info "Restart: sudo systemctl restart ${SERVICE_NAME}"
echo "================================"
