# 🐳 Docker Compose - SAMRS Project

Docker Compose setup untuk menjalankan seluruh stack SAMRS.

## 📦 Services

| Service | Port | Description |
|---------|------|-------------|
| postgres | 5432 | PostgreSQL 15 Database |
| backend | 8090 | Go API Server |
| bff | 3000 | Hono Proxy Server |
| frontend | 8000 (prod) / 5173 (dev) | React Application |
| pgadmin | 8081 | Database Management (optional) |
| prometheus | 9090 | Metrics Collection (optional) |
| grafana | 3001 | Monitoring Dashboard (optional) |

## 🚀 Quick Start

### 1. Setup Environment

```bash
# Copy environment file
cp .env.example .env

# Edit .env and change JWT_SECRET and other secrets
nano .env
```

### 2. Run All Services

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Check status
docker-compose ps
```

### 3. Initialize Database

```bash
# Run migrations and seed
docker-compose exec backend ./main seed

# Or manually
docker-compose exec backend sh
./main seed
exit
```

### 4. Access Applications

- **Frontend (Dev)**: http://localhost:5173 (with hot reload)
- **Frontend (Prod)**: http://localhost:8000 (static build)
- **Backend API**: http://localhost:8090
- **BFF**: http://localhost:3000
- **pgAdmin**: http://localhost:8081 (admin@samrs.local / admin)
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001 (admin / admin)

**Default Login:**
- Username: `admin`
- Password: `password123`
- Tenant: `rs-pusat`

## 🛠️ Development Commands

### Start Services

```bash
# Start all services (production mode)
docker-compose up -d

# Start in development mode (with hot reload)
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
# Or using Makefile
make up-dev

# Start specific service
docker-compose up -d backend

# Start with tools (pgAdmin)
docker-compose --profile tools up -d

# Start with monitoring (Prometheus + Grafana)
docker-compose --profile monitoring up -d

# Start everything
docker-compose --profile tools --profile monitoring up -d
```

### Development vs Production Mode

**Development Mode** (Recommended for development):
- Uses Vite dev server with hot reload
- Port: 5173
- Changes apply instantly without rebuild
- Command: `make up-dev`

**Production Mode**:
- Uses Nginx with static build
- Port: 8000
- Requires rebuild for every change
- Command: `make up`

### Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data!)
docker-compose down -v
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Rebuild Services

```bash
# Rebuild all (production)
docker-compose build

# Rebuild all (development)
docker-compose -f docker-compose.yml -f docker-compose.dev.yml build

# Rebuild specific service
docker-compose build backend

# Rebuild and restart (production)
docker-compose up -d --build

# Rebuild and restart (development)
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
# Or using Makefile
make up-dev
```

### Execute Commands

```bash
# Backend shell
docker-compose exec backend sh

# Run Go commands
docker-compose exec backend go test ./...

# Database shell
docker-compose exec postgres psql -U samrs_user -d samrs_db
```

## 🔧 Configuration

### Environment Variables

Edit `.env` file:

```env
JWT_SECRET=your-secret-key-min-32-chars
SEED_ADMIN_USERNAME=admin
SEED_ADMIN_PASSWORD=password123
METRICS_TOKEN=your-metrics-token
```

### Ports

Edit `docker-compose.yml` to change ports:

```yaml
services:
  backend:
    ports:
      - "8080:8080"  # Change left side: "9000:8080"
```

## 📊 Monitoring

### Enable Monitoring Stack

```bash
docker-compose --profile monitoring up -d
```

Access:
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001

### Metrics Endpoint

```bash
# Access metrics (requires token)
curl -H "Metrics-Token: secret-metrics-token" http://localhost:8080/metrics
```

## 🗄️ Database Management

### Using pgAdmin

```bash
# Start pgAdmin
docker-compose --profile tools up -d pgadmin

# Access: http://localhost:8081
# Login: admin@samrs.local / admin
```

**Add Server in pgAdmin:**
- Host: `postgres`
- Port: `5432`
- Database: `samrs_db`
- Username: `samrs_user`
- Password: `samrs_password`

### Direct Database Access

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U samrs_user -d samrs_db

# Backup database
docker-compose exec postgres pg_dump -U samrs_user samrs_db > backup.sql

# Restore database
docker-compose exec -T postgres psql -U samrs_user -d samrs_db < backup.sql
```

## 🧪 Testing

```bash
# Run backend tests
docker-compose exec backend go test ./...

# Run with coverage
docker-compose exec backend go test -cover ./...
```

## 🐛 Troubleshooting

### Backend won't start

```bash
# Check logs
docker-compose logs backend

# Check database connection
docker-compose exec backend sh
ping postgres
```

### Database connection failed

```bash
# Check if postgres is healthy
docker-compose ps postgres

# Restart postgres
docker-compose restart postgres

# Check postgres logs
docker-compose logs postgres
```

### Frontend can't connect to backend

```bash
# Check BFF logs
docker-compose logs bff

# Check backend is running
curl http://localhost:8090/api/v1/health

# Check network
docker-compose exec frontend ping backend
```

### Frontend hot reload not working

```bash
# Make sure you're using dev mode
make up-dev

# Check if Vite dev server is running
docker-compose logs frontend
# Should see: "VITE ready in XXXms"

# If using production mode, rebuild is required
docker-compose build frontend
docker-compose up -d frontend
```

### Port already in use

```bash
# Find process using port
lsof -i :5173
lsof -i :8090

# Kill process or change port in docker-compose.yml
```

### Changes not reflected after login

```bash
# Clear browser localStorage
# In browser console:
localStorage.clear()
location.reload()

# Then login again
```

## 🔄 Update & Rebuild

```bash
# Pull latest code
git pull

# Rebuild and restart (production)
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# Rebuild and restart (development with hot reload)
docker-compose down
docker-compose -f docker-compose.yml -f docker-compose.dev.yml build --no-cache
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
# Or using Makefile
make up-dev

# Run migrations if needed
docker-compose exec backend ./main migrate
```

## 🎯 Makefile Commands

Quick commands using Makefile:

```bash
# Show all available commands
make help

# Start services (production)
make up

# Start services (development with hot reload)
make up-dev

# Stop services
make down

# View logs
make logs
make logs-backend
make logs-frontend

# Restart services
make restart
make restart-backend
make restart-frontend

# Rebuild everything
make rebuild
```

## 🧹 Cleanup

```bash
# Stop and remove containers
docker-compose down

# Remove volumes (WARNING: deletes all data!)
docker-compose down -v

# Remove images
docker-compose down --rmi all

# Full cleanup
docker-compose down -v --rmi all --remove-orphans
```

## 📝 Notes

- **Production**: Change all secrets in `.env`
- **Volumes**: Data persists in Docker volumes
- **Logs**: Available in `samrs-backend/logs/`
- **Uploads**: Stored in `samrs-backend/uploads/`
- **Network**: All services in `samrs_network`

## 🔐 Security Checklist

Before deploying to production:

- [ ] Change `JWT_SECRET` to strong random string (min 32 chars)
- [ ] Change `SEED_ADMIN_PASSWORD` to strong password
- [ ] Change `METRICS_TOKEN` to secure token
- [ ] Change pgAdmin password
- [ ] Change Grafana password
- [ ] Enable HTTPS/TLS
- [ ] Configure firewall rules
- [ ] Set up backup strategy
- [ ] Configure log rotation
- [ ] Review rate limiting settings

## 📚 Additional Resources

- [Main Documentation](./DOKUMENTASI_PROJECT_SAMRS.md)
- [Backend Deep Dive](./docs/BACKEND_DEEP_DIVE.md)
- [Frontend Deep Dive](./docs/FRONTEND_DEEP_DIVE.md)
- [Development Workflow](./docs/DEVELOPMENT_WORKFLOW.md)

---

**Happy Dockerizing! 🐳**
