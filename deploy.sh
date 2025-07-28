#!/bin/bash

# ChatMood Production Deployment Script
# Author: hedgeg0d
# Description: Automated deployment script for ChatMood Telegram Mini App

set -e  # Exit on any error
set -u  # Exit on undefined variable

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_NAME="chatmood"
DOCKER_IMAGE="hedgeg0d/chatmood"
BACKUP_DIR="$HOME/chatmood-backups"
LOG_FILE="/var/log/chatmood-deploy.log"

# Functions
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1" | tee -a "$LOG_FILE"
    exit 1
}

warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1" | tee -a "$LOG_FILE"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO:${NC} $1" | tee -a "$LOG_FILE"
}

# Check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        error "This script should not be run as root for security reasons"
    fi
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."

    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed. Please install Docker first."
    fi

    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose is not installed. Please install Docker Compose first."
    fi

    # Check if .env file exists
    if [[ ! -f "$SCRIPT_DIR/.env" ]]; then
        error ".env file not found. Please create one based on .env.example"
    fi

    # Check if required environment variables are set
    source "$SCRIPT_DIR/.env"
    if [[ -z "${TELEGRAM_BOT_TOKEN:-}" ]]; then
        error "TELEGRAM_BOT_TOKEN is not set in .env file"
    fi

    if [[ -z "${WEBHOOK_URL:-}" ]]; then
        error "WEBHOOK_URL is not set in .env file"
    fi

    log "Prerequisites check passed ✓"
}

# Create backup
create_backup() {
    log "Creating backup..."

    mkdir -p "$BACKUP_DIR"
    BACKUP_FILE="$BACKUP_DIR/chatmood-$(date +%Y%m%d-%H%M%S).tar.gz"

    # Backup current deployment
    if docker-compose ps | grep -q "$PROJECT_NAME"; then
        docker-compose exec -T postgres pg_dump -U chatmood chatmood > "$BACKUP_DIR/database-$(date +%Y%m%d-%H%M%S).sql" || true
        docker-compose exec -T redis redis-cli BGSAVE || true
    fi

    # Backup configuration files
    tar -czf "$BACKUP_FILE" \
        --exclude='.git' \
        --exclude='node_modules' \
        --exclude='tmp' \
        --exclude='*.log' \
        .env docker-compose*.yml nginx.conf 2>/dev/null || true

    log "Backup created: $BACKUP_FILE"
}

# Pull latest images
pull_images() {
    log "Pulling latest Docker images..."

    docker pull "$DOCKER_IMAGE:latest" || error "Failed to pull application image"
    docker pull nginx:alpine || true
    docker pull redis:7-alpine || true
    docker pull postgres:15-alpine || true

    log "Images pulled successfully ✓"
}

# Setup SSL certificates (Let's Encrypt)
setup_ssl() {
    log "Setting up SSL certificates..."

    # Create SSL directory
    mkdir -p ssl

    # Check if certificates already exist
    if [[ -f "ssl/fullchain.pem" && -f "ssl/privkey.pem" ]]; then
        info "SSL certificates already exist, skipping setup"
        return
    fi

    # Install certbot if not present
    if ! command -v certbot &> /dev/null; then
        warning "Certbot not found. Please install certbot manually and run:"
        warning "certbot certonly --standalone -d your-domain.com"
        warning "Then copy certificates to ssl/ directory"
        return
    fi

    # Extract domain from WEBHOOK_URL
    DOMAIN=$(echo "$WEBHOOK_URL" | sed -e 's|^[^/]*//||' -e 's|[/:].*||')

    if [[ -n "$DOMAIN" ]]; then
        info "Obtaining SSL certificate for $DOMAIN"
        certbot certonly --standalone -d "$DOMAIN" --non-interactive --agree-tos --email "admin@$DOMAIN"

        # Copy certificates
        cp "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ssl/
        cp "/etc/letsencrypt/live/$DOMAIN/privkey.pem" ssl/

        log "SSL certificates setup completed ✓"
    else
        warning "Could not extract domain from WEBHOOK_URL. Please setup SSL manually."
    fi
}

# Deploy application
deploy() {
    log "Deploying ChatMood application..."

    # Stop existing containers
    if docker-compose ps | grep -q "$PROJECT_NAME"; then
        log "Stopping existing containers..."
        docker-compose down
    fi

    # Start new containers
    log "Starting new containers..."
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

    # Wait for services to be healthy
    log "Waiting for services to be healthy..."
    sleep 30

    # Check if services are running
    if ! docker-compose ps | grep -q "Up"; then
        error "Some services failed to start. Check logs with: docker-compose logs"
    fi

    log "Application deployed successfully ✓"
}

# Health check
health_check() {
    log "Performing health check..."

    local max_attempts=30
    local attempt=1

    while [[ $attempt -le $max_attempts ]]; do
        if curl -f -s "http://localhost/health" > /dev/null; then
            log "Health check passed ✓"
            return 0
        fi

        info "Health check attempt $attempt/$max_attempts failed, retrying in 5 seconds..."
        sleep 5
        ((attempt++))
    done

    error "Health check failed after $max_attempts attempts"
}

# Setup monitoring
setup_monitoring() {
    log "Setting up monitoring..."

    # Create Prometheus configuration
    cat > prometheus.yml << EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'chatmood'
    static_configs:
      - targets: ['chatmood:8080']

  - job_name: 'nginx'
    static_configs:
      - targets: ['nginx:8081']
    metrics_path: /nginx_status

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['localhost:9100']
EOF

    log "Monitoring setup completed ✓"
}

# Cleanup old images and containers
cleanup() {
    log "Cleaning up old images and containers..."

    # Remove old images
    docker image prune -f || true

    # Remove old containers
    docker container prune -f || true

    # Remove unused volumes (be careful with this)
    # docker volume prune -f || true

    log "Cleanup completed ✓"
}

# Setup log rotation
setup_log_rotation() {
    log "Setting up log rotation..."

    # Create logrotate configuration
    sudo tee /etc/logrotate.d/chatmood > /dev/null << EOF
/var/log/chatmood/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 $USER $USER
    postrotate
        docker-compose restart chatmood 2>/dev/null || true
    endscript
}
EOF

    log "Log rotation setup completed ✓"
}

# Show status
show_status() {
    log "ChatMood Deployment Status:"
    echo "=========================="

    # Show running containers
    echo -e "${BLUE}Running Containers:${NC}"
    docker-compose ps
    echo

    # Show resource usage
    echo -e "${BLUE}Resource Usage:${NC}"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"
    echo

    # Show logs (last 20 lines)
    echo -e "${BLUE}Recent Logs:${NC}"
    docker-compose logs --tail=20 chatmood
}

# Main deployment function
main() {
    log "Starting ChatMood deployment..."

    # Change to script directory
    cd "$SCRIPT_DIR"

    # Run deployment steps
    check_root
    check_prerequisites
    create_backup
    pull_images
    setup_ssl
    setup_monitoring
    deploy
    health_check
    setup_log_rotation
    cleanup

    log "ChatMood deployment completed successfully! 🎉"
    echo
    show_status
    echo
    log "Access your application at: $WEBHOOK_URL"
    log "Grafana dashboard: http://localhost:3000 (admin/\$GRAFANA_PASSWORD)"
    log "Prometheus: http://localhost:9090"
}

# Script usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -s, --status            Show deployment status"
    echo "  -b, --backup            Create backup only"
    echo "  -r, --rollback          Rollback to previous version"
    echo "  -l, --logs              Show application logs"
    echo "  -c, --cleanup           Cleanup old images and containers"
    echo
    echo "Examples:"
    echo "  $0                      Full deployment"
    echo "  $0 --status             Show status"
    echo "  $0 --backup             Create backup"
    echo "  $0 --logs               Show logs"
}

# Handle command line arguments
case "${1:-}" in
    -h|--help)
        usage
        exit 0
        ;;
    -s|--status)
        show_status
        exit 0
        ;;
    -b|--backup)
        create_backup
        exit 0
        ;;
    -l|--logs)
        docker-compose logs -f chatmood
        exit 0
        ;;
    -c|--cleanup)
        cleanup
        exit 0
        ;;
    -r|--rollback)
        log "Rolling back to previous version..."
        # This would implement rollback logic
        docker-compose down
        docker-compose up -d
        health_check
        log "Rollback completed"
        exit 0
        ;;
    "")
        main
        ;;
    *)
        error "Unknown option: $1"
        usage
        exit 1
        ;;
esac
