#!/bin/bash

set -e

# ============================================
# SAMRS - Sistem Aset Manajemen Rumah Sakit
# Management Script
# ============================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

dc() {
    docker compose "$@"
}

print_header() {
    echo ""
    echo -e "${CYAN}╔══════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║         SAMRS Management Script          ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════╝${NC}"
    echo ""
}

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_error()   { echo -e "${RED}✗ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_info()    { echo -e "${BLUE}ℹ $1${NC}"; }

check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"; exit 1
    fi
    if ! docker info &> /dev/null; then
        print_error "Docker daemon is not running"; exit 1
    fi
}

ensure_env() {
    if [ ! -f .env ]; then
        print_warning ".env file not found, creating from .env.example..."
        cp .env.example .env
        print_success ".env created"
    fi
}

show_login_info() {
    echo ""
    echo -e "${CYAN}══════════════════════════════════════════${NC}"
    echo -e "${GREEN}  SAMRS is running!${NC}"
    echo -e "${CYAN}══════════════════════════════════════════${NC}"
    echo ""
    echo -e "  ${BLUE}Access URLs:${NC}"
    echo "  ┌─────────────────┬──────────────────────┐"
    echo "  │ Frontend (Prod) │ http://localhost:8000│"
    echo "  │ Frontend (Dev)  │ http://localhost:5173│"
    echo "  │ Backend API     │ http://localhost:8090│"
    echo "  │ BFF Proxy       │ http://localhost:3000│"
    echo "  │ pgAdmin         │ http://localhost:8081│"
    echo "  │ Prometheus      │ http://localhost:9090│"
    echo "  │ Grafana         │ http://localhost:3001│"
    echo "  └─────────────────┴──────────────────────┘"
    echo ""
    echo -e "  ${YELLOW}Default Login:${NC}"
    echo "  ┌──────────┬────────────────┐"
    echo "  │ Username │ admin          │"
    echo "  │ Password │ password123    │"
    echo "  │ Tenant   │ rs-pusat       │"
    echo "  └──────────┴────────────────┘"
    echo ""
}

# ============================================
# Commands
# ============================================

cmd_start() {
    local mode="${1:-prod}"
    check_docker
    ensure_env
    print_header
    print_info "Starting SAMRS ($mode mode)..."

    case "$mode" in
        dev)
            dc -f docker-compose.yml -f docker-compose.dev.yml up -d
            ;;
        tools)
            dc --profile tools up -d
            ;;
        monitoring)
            dc --profile monitoring up -d
            ;;
        all)
            dc --profile tools --profile monitoring up -d
            ;;
        *)
            dc up -d
            ;;
    esac

    print_success "Services started"
    dc ps
    show_login_info
}

cmd_stop() {
    check_docker
    print_info "Stopping SAMRS..."
    dc down
    print_success "All services stopped"
}

cmd_restart() {
    check_docker
    local service="$1"
    if [ -n "$service" ]; then
        print_info "Restarting $service..."
        dc restart "$service"
        print_success "$service restarted"
    else
        print_info "Restarting all services..."
        dc restart
        print_success "All services restarted"
    fi
}

cmd_status() {
    check_docker
    print_info "SAMRS Services Status:"
    echo ""
    dc ps
    echo ""
    echo -e "${BLUE}Health Check:${NC}"
    curl -s http://localhost:8090/ping > /dev/null 2>&1 && print_success "Backend (8090): UP" || print_error "Backend (8090): DOWN"
    curl -s http://localhost:8000 > /dev/null 2>&1 && print_success "Frontend (8000): UP" || print_error "Frontend (8000): DOWN"
    curl -s http://localhost:3000 > /dev/null 2>&1 && print_success "BFF (3000): UP" || print_error "BFF (3000): DOWN"
    echo ""
}

cmd_logs() {
    check_docker
    local service="$1"
    if [ -n "$service" ]; then
        dc logs -f "$service"
    else
        dc logs -f
    fi
}

cmd_build() {
    check_docker
    print_info "Building all services..."
    dc build
    print_success "Build completed"
}

cmd_rebuild() {
    check_docker
    print_info "Rebuilding all services (no cache)..."
    dc down
    dc build --no-cache
    dc up -d
    print_success "Rebuild completed"
    dc ps
    show_login_info
}

cmd_seed() {
    check_docker
    print_info "Running database seeder..."
    dc exec backend ./seed
    print_success "Database seeded"
}

cmd_test() {
    check_docker
    print_info "Running backend tests..."
    dc exec backend go test ./...
}

cmd_test_coverage() {
    check_docker
    print_info "Running backend tests with coverage..."
    dc exec backend go test -cover ./...
}

cmd_shell() {
    check_docker
    local target="${1:-backend}"
    case "$target" in
        backend|be)
            dc exec backend sh
            ;;
        db|database)
            dc exec postgres psql -U samrs_user -d samrs_db
            ;;
        *)
            print_error "Unknown shell target: $target"
            print_info "Available: backend, db"
            exit 1
            ;;
    esac
}

cmd_backup() {
    check_docker
    local filename="backup_$(date +%Y%m%d_%H%M%S).sql"
    print_info "Backing up database to $filename..."
    dc exec -T postgres pg_dump -U samrs_user samrs_db > "$filename"
    print_success "Database backed up to $filename"
}

cmd_restore() {
    check_docker
    local file="$1"
    if [ -z "$file" ]; then
        print_error "Usage: ./samrs.sh restore <backup_file.sql>"
        exit 1
    fi
    if [ ! -f "$file" ]; then
        print_error "File not found: $file"
        exit 1
    fi
    print_warning "This will overwrite the current database!"
    read -p "Are you sure? (y/N): " confirm
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        print_info "Restore cancelled"; exit 0
    fi
    print_info "Restoring database from $file..."
    dc exec -T postgres psql -U samrs_user -d samrs_db < "$file"
    print_success "Database restored from $file"
}

cmd_clean() {
    check_docker
    print_warning "This will remove ALL containers, volumes, and images!"
    read -p "Are you sure? (y/N): " confirm
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        print_info "Clean cancelled"; exit 0
    fi
    print_info "Cleaning up..."
    dc down -v --rmi all --remove-orphans
    print_success "All resources cleaned"
}

cmd_help() {
    print_header
    echo -e "${GREEN}Usage:${NC} ./samrs.sh <command> [options]"
    echo ""
    echo -e "${YELLOW}Start Commands:${NC}"
    echo "  start              Start services (production mode)"
    echo "  start:dev          Start with hot reload (development)"
    echo "  start:tools        Start with pgAdmin"
    echo "  start:monitoring   Start with Prometheus & Grafana"
    echo "  start:all          Start everything"
    echo ""
    echo -e "${YELLOW}Stop & Restart:${NC}"
    echo "  stop               Stop all services"
    echo "  restart            Restart all services"
    echo "  restart <service>  Restart specific service"
    echo ""
    echo -e "${YELLOW}Status & Logs:${NC}"
    echo "  status             Show services status & health"
    echo "  logs               Show all logs"
    echo "  logs <service>     Show specific service logs"
    echo ""
    echo -e "${YELLOW}Build:${NC}"
    echo "  build              Build all Docker images"
    echo "  rebuild            Rebuild (no cache) + restart"
    echo ""
    echo -e "${YELLOW}Database:${NC}"
    echo "  seed               Run database seeder"
    echo "  db                 Open PostgreSQL shell"
    echo "  backup             Backup database"
    echo "  restore <file>     Restore database from file"
    echo ""
    echo -e "${YELLOW}Testing:${NC}"
    echo "  test               Run backend tests"
    echo "  test:cov           Run with coverage"
    echo ""
    echo -e "${YELLOW}Shell:${NC}"
    echo "  shell              Open backend shell"
    echo ""
    echo -e "${YELLOW}Utility:${NC}"
    echo "  clean              Remove all containers/volumes/images"
    echo "  help               Show this help"
    echo ""
}

# ============================================
# Main Router
# ============================================

case "${1:-help}" in
    start)              cmd_start prod ;;
    start:dev)          cmd_start dev ;;
    start:tools)        cmd_start tools ;;
    start:monitoring)   cmd_start monitoring ;;
    start:all)          cmd_start all ;;
    stop)               cmd_stop ;;
    restart)            cmd_restart "$2" ;;
    status|st)          cmd_status ;;
    logs|log)           cmd_logs "$2" ;;
    build)              cmd_build ;;
    rebuild)            cmd_rebuild ;;
    seed)               cmd_seed ;;
    test)               cmd_test ;;
    test:cov)           cmd_test_coverage ;;
    shell|sh)           cmd_shell "$2" ;;
    db|database)        cmd_shell db ;;
    backup)             cmd_backup ;;
    restore)            cmd_restore "$2" ;;
    clean)              cmd_clean ;;
    help|--help|-h)     cmd_help ;;
    *)
        print_error "Unknown command: $1"
        cmd_help
        exit 1
        ;;
esac
