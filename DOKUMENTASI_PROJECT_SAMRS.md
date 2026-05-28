# 📘 DOKUMENTASI PROJECT SAMRS (Sistem Aset Manajemen Rumah Sakit)

> **Dokumentasi Teknis Lengkap untuk Developer**  
> Versi: 2.0 (Phase 2 - Foundation Hardening)  
> Terakhir Diperbarui: 28 Februari 2026

---

## 📋 DAFTAR ISI

1. [Overview Project](#1-overview-project)
2. [Arsitektur Sistem](#2-arsitektur-sistem)
3. [Tech Stack](#3-tech-stack)
4. [Struktur Project](#4-struktur-project)
5. [Backend Deep Dive](./docs/BACKEND_DEEP_DIVE.md)
6. [Frontend Deep Dive](./docs/FRONTEND_DEEP_DIVE.md)
7. [Database Schema](./docs/DATABASE_SCHEMA.md)
8. [Authentication & Authorization](./docs/AUTH_RBAC.md)
9. [Multi-Tenancy Implementation](./docs/MULTI_TENANCY.md)
10. [Audit Trail System](./docs/AUDIT_TRAIL.md)
11. [API Documentation](./docs/API_DOCUMENTATION.md)
12. [Development Workflow](./docs/DEVELOPMENT_WORKFLOW.md)
13. [Testing Strategy](./docs/TESTING_STRATEGY.md)
14. [Deployment](./docs/DEPLOYMENT.md)
15. [Troubleshooting](./docs/TROUBLESHOOTING.md)

---

## 1. OVERVIEW PROJECT

### 1.1 Apa itu SAMRS?

SAMRS (Sistem Aset Manajemen Rumah Sakit) adalah platform **SaaS Multi-Tenant** untuk manajemen aset medis dan non-medis di rumah sakit. Sistem ini dibangun dengan fokus pada:

- **Multi-Tenancy**: Isolasi data total antar rumah sakit
- **RBAC (Role-Based Access Control)**: Kontrol akses berbasis permission yang fleksibel
- **Asset Lifecycle Management**: Tracking lengkap dari registrasi hingga disposal
- **Audit Trail**: Pencatatan setiap perubahan data untuk compliance
- **Security First**: Tenant isolation, authentication, dan authorization yang ketat

### 1.2 Business Requirements

**Target Users:**
1. **Super Admin (Global)**: Mengelola tenant, monitoring sistem
2. **Tenant Admin (RS)**: Mengelola user, role, dan konfigurasi RS
3. **Technician**: Maintenance, update status aset
4. **Staff/Nurse**: Melaporkan kerusakan, melihat ketersediaan aset

**Core Features (MVP):**
- ✅ Multi-tenant dengan isolasi data
- ✅ Authentication & Authorization (JWT + RBAC)
- ✅ Master Data (Room, Bed, Category, Vendor, Brand, Model, Status)
- ✅ Asset Management (CRUD, QR Code, Timeline)
- ✅ Asset Mutation (Perpindahan lokasi)
- ✅ Stock Opname (Inventarisasi berkala)
- ✅ Complaint Management (Pelaporan kerusakan)
- ✅ Maintenance Scheduling (Jadwal pemeliharaan)
- ✅ Document Management (Upload/Download dokumen)
- ✅ Reporting (Export CSV/XLSX/PDF)
- ✅ Audit Trail (Tracking semua perubahan)

### 1.3 Project Status

**Current Phase**: Phase 2 - Foundation Hardening

**Completed:**
- ✅ MVP Features (All core modules)
- ✅ Multi-tenant architecture
- ✅ RBAC implementation
- ✅ Audit trail system
- ✅ API documentation (Swagger)
- ✅ Monitoring (Prometheus + Grafana)
- ✅ Structured logging
- ✅ Rate limiting
- ✅ Metrics endpoint

**In Progress (Phase 2):**
- 🔄 Auth & RBAC hardening
- 🔄 Multi-tenancy safety net
- 🔄 Data integrity & consistency
- 🔄 Performance baseline
- 🔄 API contract freeze

---

## 2. ARSITEKTUR SISTEM

### 2.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Browser    │  │    Mobile    │  │   Desktop    │      │
│  │   (React)    │  │    (Future)  │  │   (Future)   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    FRONTEND LAYER (BFF)                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  React 19 + Vite + TypeScript                        │   │
│  │  - Redux Toolkit + RTK Query                         │   │
│  │  - TanStack Router + Table                           │   │
│  │  - Shadcn/UI + Tailwind CSS                          │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  BFF (Backend for Frontend) - Hono.js                │   │
│  │  - CORS handling                                     │   │
│  │  - Request proxying                                  │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                     BACKEND LAYER (API)                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Go 1.25 + Gin Framework                             │   │
│  │  ┌────────────────────────────────────────────────┐  │   │
│  │  │  Middleware Layer                              │  │   │
│  │  │  - Request ID                                  │  │   │
│  │  │  - Structured Logger                           │  │   │
│  │  │  - Rate Limiter                                │  │   │
│  │  │  - Auth (JWT)                                  │  │   │
│  │  │  - RBAC (Permission Check)                     │  │   │
│  │  │  - Tenant Scope                                │  │   │
│  │  │  - Metrics (Prometheus)                        │  │   │
│  │  └────────────────────────────────────────────────┘  │   │
│  │  ┌────────────────────────────────────────────────┐  │   │
│  │  │  Handler Layer (HTTP Controllers)              │  │   │
│  │  └────────────────────────────────────────────────┘  │   │
│  │  ┌────────────────────────────────────────────────┐  │   │
│  │  │  Usecase Layer (Business Logic)                │  │   │
│  │  └────────────────────────────────────────────────┘  │   │
│  │  ┌────────────────────────────────────────────────┐  │   │
│  │  │  Repository Layer (Data Access)                │  │   │
│  │  └────────────────────────────────────────────────┘  │   │
│  │  ┌────────────────────────────────────────────────┐  │   │
│  │  │  Domain Layer (Entities & Models)              │  │   │
│  │  └────────────────────────────────────────────────┘  │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    DATABASE LAYER                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  PostgreSQL 15                                       │   │
│  │  - Shared Schema Strategy                            │   │
│  │  - Tenant ID in every table                          │   │
│  │  - GORM ORM                                          │   │
│  │  - Audit Trail Hooks                                 │   │
│  │  - Tenant Scope Guard                                │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  MONITORING & LOGGING                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Prometheus  │  │   Grafana    │  │  Lumberjack  │      │
│  │  (Metrics)   │  │ (Dashboard)  │  │ (Log Rotate) │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Clean Architecture Pattern

Project ini menggunakan **Clean Architecture** dengan layer separation:

1. **Domain Layer**: Entity & business rules (tidak depend ke layer lain)
2. **Repository Layer**: Data access abstraction
3. **Usecase Layer**: Business logic & orchestration
4. **Delivery Layer**: HTTP handlers, middleware, routing

**Dependency Rule**: Inner layer tidak boleh depend ke outer layer.

```
Domain ← Repository ← Usecase ← Delivery
```

### 2.3 Request Flow

```
1. Client Request
   ↓
2. BFF (Hono) - CORS & Proxy
   ↓
3. Gin Router
   ↓
4. Middleware Chain:
   - Request ID
   - Logger
   - Rate Limiter
   - Auth (JWT validation)
   - Tenant Scope
   - RBAC (Permission check)
   ↓
5. Handler (HTTP Controller)
   - Parse request
   - Validate input
   - Call usecase
   ↓
6. Usecase (Business Logic)
   - Validate business rules
   - Call repository
   - Orchestrate operations
   ↓
7. Repository (Data Access)
   - GORM queries
   - Tenant scoped
   - Audit trail hooks
   ↓
8. Database (PostgreSQL)
   ↓
9. Response back through layers
   ↓
10. Client receives JSON response
```

---

## 3. TECH STACK

### 3.1 Backend Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Go | 1.25.5 | Backend language |
| Framework | Gin | 1.11.0 | HTTP web framework |
| ORM | GORM | 1.31.1 | Database ORM |
| Database | PostgreSQL | 15 | Primary database |
| Auth | JWT | 5.3.0 | Authentication |
| Validation | go-playground/validator | 10.30.1 | Input validation |
| Logger | Zap | 1.27.0 | Structured logging |
| Log Rotation | Lumberjack | 2.2.1 | Log file rotation |
| Metrics | Prometheus | 1.23.2 | Monitoring metrics |
| QR Code | go-qrcode | - | QR code generation |
| Excel | excelize | 2.10.0 | Excel export |
| PDF | gofpdf | 1.16.2 | PDF generation |
| Rate Limit | golang.org/x/time | 0.12.0 | Rate limiting |
| Testing | go.uber.org/mock | 0.6.0 | Mock generation |

### 3.2 Frontend Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | TypeScript | 5.9.3 | Type-safe JavaScript |
| Framework | React | 19.2.0 | UI framework |
| Build Tool | Vite (Rolldown) | 7.2.5 | Fast build tool |
| State Management | Redux Toolkit | 2.11.2 | Global state |
| API Client | RTK Query | - | Data fetching & caching |
| Router | TanStack Router | 1.147.3 | Type-safe routing |
| Table | TanStack Table | 8.21.3 | Headless table |
| Virtual | TanStack Virtual | 3.13.18 | Virtualization |
| UI Components | Shadcn/UI | - | Accessible components |
| Styling | Tailwind CSS | 4.1.18 | Utility-first CSS |
| Forms | React Hook Form | 7.71.0 | Form management |
| Validation | Zod | 4.3.5 | Schema validation |
| Date | date-fns | 4.1.0 | Date manipulation |
| Icons | Lucide React | 0.562.0 | Icon library |
| Toast | Sonner | 2.0.7 | Toast notifications |
| DnD | dnd-kit | 6.3.1 | Drag and drop |
| BFF | Hono | 4.11.4 | Edge framework |

### 3.3 DevOps & Tools

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Container | Docker | Containerization |
| Database Admin | pgAdmin | PostgreSQL GUI |
| Monitoring | Prometheus + Grafana | Metrics & dashboards |
| API Testing | REST Client (VS Code) | API testing |
| API Docs | Swagger/OpenAPI | API documentation |
| Version Control | Git | Source control |

---

## 4. STRUKTUR PROJECT

### 4.1 Root Structure

```
samrs-project/
├── samrs-backend/          # Backend Go application
├── samrs-frontend/         # Frontend React application
├── dokumentasi-be/         # Backend documentation
├── dokumentasi-fe/         # Frontend documentation
├── monitoring/             # Prometheus & Grafana configs
└── docs/                   # Project documentation
    ├── BACKEND_DEEP_DIVE.md
    ├── FRONTEND_DEEP_DIVE.md
    ├── DATABASE_SCHEMA.md
    ├── AUTH_RBAC.md
    ├── MULTI_TENANCY.md
    ├── AUDIT_TRAIL.md
    ├── API_DOCUMENTATION.md
    ├── DEVELOPMENT_WORKFLOW.md
    ├── TESTING_STRATEGY.md
    ├── DEPLOYMENT.md
    └── TROUBLESHOOTING.md
```

### 4.2 Backend Structure (Simplified)

```
samrs-backend/
├── cmd/
│   ├── api/main.go         # Entry point
│   └── seed/main.go        # Database seeder
├── internal/
│   ├── app/                # DI & routing
│   ├── config/             # DB, logger, hooks
│   ├── domain/             # 20 entities
│   ├── repository/         # Data access (40+ repos)
│   ├── usecase/            # Business logic (40+ usecases)
│   ├── delivery/http/      # Handlers & middleware
│   ├── seed/               # Seeder logic
│   └── test/mocks/         # Generated mocks
├── pkg/util/               # Utilities
├── uploads/                # File storage
├── logs/                   # Application logs
├── .env                    # Environment config
├── docker-compose.yml      # Docker setup
└── api_test.http           # API tests
```

**Stats:**
- 182 Go files
- ~31,733 lines of code
- 20 domain models
- 40+ repositories
- 40+ use cases
- 20+ handlers

### 4.3 Frontend Structure (Simplified)

```
samrs-frontend/
├── apps/
│   ├── frontend/           # Main React app
│   │   └── src/
│   │       ├── components/ # UI components
│   │       ├── constants/  # Config & permissions
│   │       ├── hooks/      # Custom hooks
│   │       ├── modules/    # Feature modules
│   │       │   ├── auth/
│   │       │   ├── asset/
│   │       │   ├── mutation/
│   │       │   ├── dashboard/
│   │       │   ├── forms/
│   │       │   └── ...
│   │       ├── routes/     # Routing
│   │       ├── store/      # Redux store
│   │       └── main.tsx    # Entry point
│   └── bff/                # Hono proxy
│       └── src/index.ts
└── packages/               # Shared (future)
```

**Stats:**
- 112 TypeScript/TSX files
- 7+ feature modules
- 40+ UI components
- Modular by feature

---

## 5. QUICK START

### 5.1 Prerequisites

- Go 1.25+
- Node.js 18+ (with Bun)
- PostgreSQL 15
- Docker & Docker Compose (optional)

### 5.2 Backend Setup

```bash
# 1. Clone repository
cd samrs-backend

# 2. Copy environment file
cp .env.example .env

# 3. Edit .env with your database credentials
nano .env

# 4. Start PostgreSQL (via Docker)
docker-compose up -d db

# 5. Install dependencies
go mod download

# 6. Run migrations & seed
go run cmd/seed/main.go

# 7. Start server
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

### 5.3 Frontend Setup

```bash
# 1. Navigate to frontend
cd samrs-frontend

# 2. Install dependencies
bun install

# 3. Start BFF & Frontend
bun dev
```

- BFF: `http://localhost:3000`
- Frontend: `http://localhost:5173`

### 5.4 Default Credentials

```
Username: admin
Password: password123
Tenant: rs-pusat
```

---

## 6. KEY CONCEPTS

### 6.1 Multi-Tenancy

Setiap tenant (rumah sakit) memiliki data yang terisolasi:

- Setiap tabel memiliki kolom `tenant_id`
- GORM Tenant Scope Guard memastikan query selalu filtered by tenant
- Middleware `TenantScope` inject tenant_id dari JWT ke context
- Super Admin dapat bypass tenant scope untuk management

### 6.2 RBAC (Role-Based Access Control)

Permission-based authorization:

- Setiap endpoint dilindungi oleh permission slug (e.g., `asset:create`)
- User memiliki Role, Role memiliki Permissions
- Middleware `RBACMiddleware` check permission sebelum akses endpoint
- Flexible: Tenant Admin dapat assign/revoke permissions

### 6.3 Audit Trail

Automatic logging semua perubahan data:

- GORM hooks capture CREATE, UPDATE, DELETE
- Menyimpan old_value dan new_value dalam JSON
- Mencatat user_id, tenant_id, IP address, timestamp
- Tidak bisa dimatikan (kecuali explicit flag)

### 6.4 Clean Architecture

Separation of concerns:

```
Domain (entities) 
  ↑
Repository (data access interface)
  ↑
Usecase (business logic)
  ↑
Handler (HTTP controller)
```

Inner layer tidak depend ke outer layer.

---

## 7. DEVELOPMENT GUIDELINES

### 7.1 Backend Guidelines

**Menambah Modul Baru:**

1. Buat domain entity di `internal/domain/`
2. Buat repository interface & implementation di `internal/repository/`
3. Buat usecase interface & implementation di `internal/usecase/`
4. Buat handler di `internal/delivery/http/handlers/`
5. Register routes di `internal/app/routes/`
6. Wire dependencies di `internal/app/app.go`

**Aturan Penting:**
- Semua query HARUS tenant-scoped
- Semua endpoint HARUS memiliki permission check
- Gunakan transaction untuk operasi multi-table
- Validasi input di handler, business rules di usecase
- Error handling menggunakan custom error types

### 7.2 Frontend Guidelines

**Menambah Modul Baru:**

1. Buat folder di `src/modules/[module-name]/`
2. Struktur: `actions/`, `components/`, `hooks/`, `layouts/`, `pages/`, `routes.tsx`, `types.ts`, `schemas.ts`
3. Define API di `actions/[module]Api.ts` menggunakan RTK Query
4. Define Redux slice di `actions/[module]Slice.ts` (jika perlu state)
5. Buat routes di `routes.tsx`
6. Register routes di `src/routes/index.tsx`

**Aturan Penting:**
- Gunakan TypeScript strict mode
- Semua API call via RTK Query
- Form validation menggunakan Zod + React Hook Form
- Permission check menggunakan `usePermission` hook
- Responsive design (mobile-first)

---

## 8. TESTING

### 8.1 Backend Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/usecase/asset/...
```

**Test Coverage:**
- Unit tests untuk usecase layer
- Integration tests untuk repository layer
- Mock dependencies menggunakan go.uber.org/mock

### 8.2 Frontend Testing

```bash
# Run tests (future)
bun test
```

---

## 9. MONITORING

### 9.1 Metrics Endpoint

```
GET /metrics
Header: Metrics-Token: <METRICS_TOKEN>
```

Metrics yang tersedia:
- HTTP request duration
- HTTP request count by status
- Active requests

### 9.2 Prometheus + Grafana

```bash
cd monitoring
METRICS_TOKEN=your-token docker compose -f docker-compose.monitoring.yml up -d
```

- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (admin/admin)

---

## 10. API DOCUMENTATION

### 10.1 Swagger UI

- Swagger UI: `http://localhost:8080/swagger-ui`
- Redoc: `http://localhost:8080/swagger`
- Spec: `http://localhost:8080/swagger.yaml`

### 10.2 API Test Collection

Gunakan file `samrs-backend/api_test.http` dengan REST Client extension di VS Code.

---

## 11. DEPLOYMENT

### 11.1 Production Checklist

- [ ] Set `JWT_SECRET` yang kuat
- [ ] Set `METRICS_TOKEN` yang aman
- [ ] Configure rate limiting sesuai kebutuhan
- [ ] Setup log aggregation (ELK, Loki, dll)
- [ ] Setup backup database
- [ ] Configure CORS di BFF untuk production domain
- [ ] Enable HTTPS
- [ ] Setup monitoring alerts

### 11.2 Environment Variables

Lihat `.env.example` untuk daftar lengkap environment variables.

---

## 12. TROUBLESHOOTING

### 12.1 Backend Issues

**Database connection failed:**
- Check PostgreSQL running: `docker ps`
- Check credentials di `.env`
- Check port 5432 tidak digunakan aplikasi lain

**Permission denied:**
- Check JWT token valid
- Check user memiliki permission yang sesuai
- Check role_permissions table

**Tenant data leak:**
- Check tenant_id di JWT payload
- Check middleware TenantScope aktif
- Check query menggunakan tenant scope

### 12.2 Frontend Issues

**API call failed:**
- Check BFF running di port 3000
- Check backend running di port 8080
- Check CORS configuration
- Check network tab di browser DevTools

**Permission denied:**
- Check user login dengan role yang sesuai
- Check permissions di Redux state
- Check `usePermission` hook

---

## 13. CONTRIBUTING

### 13.1 Git Workflow

```bash
# Create feature branch
git checkout -b feature/nama-fitur

# Commit changes
git add .
git commit -m "feat: deskripsi fitur"

# Push to remote
git push origin feature/nama-fitur

# Create Pull Request
```

### 13.2 Commit Convention

```
feat: new feature
fix: bug fix
docs: documentation
refactor: code refactoring
test: add tests
chore: maintenance
```

---

## 14. RESOURCES

### 14.1 Documentation

- [Backend README](./samrs-backend/README.md)
- [Frontend README](./samrs-frontend/README.md)
- [Project Status v2](./dokumentasi-be/PROJECT_STATUS_v2.md)
- [BRD](./dokumentasi-be/abstrak/01.%20BRD.md)

### 14.2 External Links

- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [React](https://react.dev/)
- [Redux Toolkit](https://redux-toolkit.js.org/)
- [TanStack Router](https://tanstack.com/router)
- [Shadcn/UI](https://ui.shadcn.com/)

---

## 15. CONTACT & SUPPORT

Untuk pertanyaan atau issue, silakan:
1. Check dokumentasi ini terlebih dahulu
2. Check existing issues di repository
3. Buat issue baru dengan detail yang jelas

---

**Happy Coding! 🚀**
