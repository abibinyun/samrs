# 📚 DOKUMENTASI PROJECT SAMRS - INDEX

Selamat datang di dokumentasi lengkap project SAMRS (Sistem Aset Manajemen Rumah Sakit).

## 📖 Dokumentasi Utama

### [DOKUMENTASI_PROJECT_SAMRS.md](./DOKUMENTASI_PROJECT_SAMRS.md)
Dokumentasi utama yang berisi overview, arsitektur, tech stack, dan quick start guide.

**Isi:**
- Overview Project
- Arsitektur Sistem
- Tech Stack
- Struktur Project
- Quick Start
- Key Concepts
- Development Guidelines

---

## 📂 Dokumentasi Detail

### 1. [Backend Deep Dive](./docs/BACKEND_DEEP_DIVE.md)
Panduan mendalam tentang backend architecture, clean architecture layers, dan implementasi detail.

**Topik:**
- Clean Architecture Layers
- Domain Layer (20 entities)
- Repository Pattern
- Usecase Layer
- Delivery Layer (Handlers & Middleware)
- Error Handling
- Response Format
- Testing
- Best Practices

### 2. [Frontend Deep Dive](./docs/FRONTEND_DEEP_DIVE.md)
Panduan lengkap frontend architecture, state management, dan component structure.

**Topik:**
- Tech Stack Overview
- Modular Structure
- State Management (Redux Toolkit + RTK Query)
- Routing (TanStack Router)
- Components & UI
- Forms & Validation
- Custom Hooks
- BFF (Backend for Frontend)
- Performance Optimization

### 3. [Database Schema](./docs/DATABASE_SCHEMA.md)
Dokumentasi lengkap database schema, relationships, dan indexing strategy.

**Topik:**
- Database Overview
- Core Tables (Tenants, Users, Roles, Permissions)
- Master Data Tables
- Asset Tables
- Operational Tables
- Document Tables
- Audit Table
- Entity Relationships
- Indexes Strategy
- Data Integrity

### 4. [Authentication & RBAC](./docs/AUTH_RBAC.md)
Panduan authentication dan authorization system.

**Topik:**
- JWT-Based Authentication
- Login Flow
- Auth Middleware
- RBAC Model
- Permission List
- Default Roles
- RBAC Middleware
- Frontend Permission Check
- Role Management
- Security Best Practices

### 5. [Development Workflow](./docs/DEVELOPMENT_WORKFLOW.md)
Panduan workflow development dari setup hingga deployment.

**Topik:**
- Setup Development Environment
- Menambah Fitur Baru (Backend & Frontend)
- Testing Strategy
- Git Workflow
- Code Review Checklist
- Debugging

---

## 🚀 Quick Links

### Untuk Developer Baru
1. Baca [DOKUMENTASI_PROJECT_SAMRS.md](./DOKUMENTASI_PROJECT_SAMRS.md) - Overview & Quick Start
2. Baca [Backend Deep Dive](./docs/BACKEND_DEEP_DIVE.md) - Memahami backend
3. Baca [Frontend Deep Dive](./docs/FRONTEND_DEEP_DIVE.md) - Memahami frontend
4. Baca [Development Workflow](./docs/DEVELOPMENT_WORKFLOW.md) - Mulai coding

### Untuk Memahami Data
1. Baca [Database Schema](./docs/DATABASE_SCHEMA.md) - Struktur database
2. Baca [Backend Deep Dive](./docs/BACKEND_DEEP_DIVE.md) - Domain entities

### Untuk Memahami Security
1. Baca [Authentication & RBAC](./docs/AUTH_RBAC.md) - Auth & permissions
2. Baca [Backend Deep Dive](./docs/BACKEND_DEEP_DIVE.md) - Middleware

---

## 📊 Project Statistics

**Backend:**
- 182 Go files
- ~31,733 lines of code
- 20 domain models
- 40+ repositories
- 40+ use cases
- 20+ handlers

**Frontend:**
- 112 TypeScript/TSX files
- 7+ feature modules
- 40+ UI components
- Modular by feature

**Database:**
- 20+ tables
- PostgreSQL 15
- Shared schema multi-tenancy

---

## 🎯 Key Features

✅ Multi-tenant SaaS
✅ JWT Authentication
✅ RBAC (Role-Based Access Control)
✅ Asset Management
✅ Asset Mutation & Timeline
✅ Stock Opname
✅ Complaint Management
✅ Maintenance Scheduling
✅ Document Management
✅ Reporting (CSV/XLSX/PDF)
✅ Audit Trail
✅ Monitoring (Prometheus + Grafana)

---

## 🛠️ Tech Stack

**Backend:**
- Go 1.25 + Gin
- GORM + PostgreSQL 15
- JWT Authentication
- Zap Logger
- Prometheus Metrics

**Frontend:**
- React 19 + TypeScript
- Vite (Rolldown)
- Redux Toolkit + RTK Query
- TanStack Router + Table
- Shadcn/UI + Tailwind CSS
- React Hook Form + Zod

---

## 📞 Support

Untuk pertanyaan atau issue:
1. Check dokumentasi ini terlebih dahulu
2. Check existing issues di repository
3. Buat issue baru dengan detail yang jelas

---

## 📝 Changelog

**2026-02-28:**
- ✅ Dokumentasi lengkap dibuat
- ✅ Backend Deep Dive
- ✅ Frontend Deep Dive
- ✅ Database Schema
- ✅ Authentication & RBAC
- ✅ Development Workflow

---

**Happy Coding! 🚀**
