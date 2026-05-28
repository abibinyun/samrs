# Database Schema - SAMRS

## 1. OVERVIEW

### 1.1 Database Info

- **DBMS**: PostgreSQL 15
- **Strategy**: Shared Schema Multi-Tenancy
- **ORM**: GORM
- **Total Tables**: 20+
- **Extensions**: uuid-ossp

### 1.2 Multi-Tenancy Strategy

**Shared Schema** dengan `tenant_id` di setiap tabel:

**Advantages:**
- Simple deployment
- Easy backup & restore
- Cost-effective
- Good for MVP

**Disadvantages:**
- Tenant scope guard critical
- Risk of data leak if not careful
- Scaling challenges (future)

---

## 2. CORE TABLES

### 2.1 Tenants

```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    address VARCHAR(255),
    tenant_type VARCHAR(50),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_status ON tenants(status);
```

### 2.2 Users

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    role_id INTEGER NOT NULL REFERENCES roles(id),
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

### 2.3 Roles

```sql
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id), -- NULL for global roles
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
CREATE INDEX idx_roles_slug ON roles(slug);
```

### 2.4 Permissions

```sql
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_permissions_slug ON permissions(slug);
```

### 2.5 Role Permissions

```sql
CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
```

---

## 3. MASTER DATA TABLES

### 3.1 Categories

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_categories_tenant_id ON categories(tenant_id);
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);
```

### 3.2 Rooms

```sql
CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    floor VARCHAR(50),
    building VARCHAR(100),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_rooms_tenant_id ON rooms(tenant_id);
CREATE INDEX idx_rooms_code ON rooms(code);
CREATE INDEX idx_rooms_deleted_at ON rooms(deleted_at);
CREATE UNIQUE INDEX idx_rooms_tenant_code ON rooms(tenant_id, code) WHERE deleted_at IS NULL;
```

### 3.3 Beds

```sql
CREATE TABLE beds (
    id SERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    room_id UUID NOT NULL REFERENCES rooms(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'available',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_beds_tenant_id ON beds(tenant_id);
CREATE INDEX idx_beds_room_id ON beds(room_id);
CREATE INDEX idx_beds_code ON beds(code);
CREATE INDEX idx_beds_status ON beds(status);
CREATE INDEX idx_beds_deleted_at ON beds(deleted_at);
CREATE UNIQUE INDEX idx_beds_tenant_code ON beds(tenant_id, code) WHERE deleted_at IS NULL;
```

### 3.4 Vendors

```sql
CREATE TABLE vendors (
    id SERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    contact VARCHAR(100),
    address TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_vendors_tenant_id ON vendors(tenant_id);
CREATE INDEX idx_vendors_code ON vendors(code);
CREATE INDEX idx_vendors_deleted_at ON vendors(deleted_at);
```

### 3.5 Asset Brands

```sql
CREATE TABLE asset_brands (
    id SERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_asset_brands_tenant_id ON asset_brands(tenant_id);
CREATE INDEX idx_asset_brands_code ON asset_brands(code);
CREATE INDEX idx_asset_brands_deleted_at ON asset_brands(deleted_at);
```

### 3.6 Asset Models

```sql
CREATE TABLE asset_models (
    id SERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    brand_id INTEGER REFERENCES asset_brands(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_asset_models_tenant_id ON asset_models(tenant_id);
CREATE INDEX idx_asset_models_brand_id ON asset_models(brand_id);
CREATE INDEX idx_asset_models_code ON asset_models(code);
CREATE INDEX idx_asset_models_deleted_at ON asset_models(deleted_at);
```

### 3.7 Asset Statuses

```sql
CREATE TABLE asset_statuses (
    id SERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(50),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_asset_statuses_tenant_id ON asset_statuses(tenant_id);
CREATE INDEX idx_asset_statuses_code ON asset_statuses(code);
CREATE INDEX idx_asset_statuses_deleted_at ON asset_statuses(deleted_at);
```

---

## 4. ASSET TABLES

### 4.1 Assets

```sql
CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    category_id INTEGER NOT NULL REFERENCES categories(id),
    room_id UUID REFERENCES rooms(id),
    bed_id INTEGER REFERENCES beds(id),
    vendor_id INTEGER REFERENCES vendors(id),
    brand_id INTEGER REFERENCES asset_brands(id),
    model_id INTEGER REFERENCES asset_models(id),
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(100),
    model VARCHAR(100),
    status VARCHAR(50) DEFAULT 'ready',
    purchase_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_assets_tenant_id ON assets(tenant_id);
CREATE INDEX idx_assets_category_id ON assets(category_id);
CREATE INDEX idx_assets_room_id ON assets(room_id);
CREATE INDEX idx_assets_bed_id ON assets(bed_id);
CREATE INDEX idx_assets_vendor_id ON assets(vendor_id);
CREATE INDEX idx_assets_brand_id ON assets(brand_id);
CREATE INDEX idx_assets_model_id ON assets(model_id);
CREATE INDEX idx_assets_code ON assets(code);
CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_assets_created_at ON assets(created_at);
CREATE INDEX idx_assets_deleted_at ON assets(deleted_at);
CREATE UNIQUE INDEX idx_assets_tenant_code ON assets(tenant_id, code) WHERE deleted_at IS NULL;

-- Composite indexes for common queries
CREATE INDEX idx_assets_tenant_status ON assets(tenant_id, status);
CREATE INDEX idx_assets_tenant_category ON assets(tenant_id, category_id);
CREATE INDEX idx_assets_tenant_room ON assets(tenant_id, room_id);
CREATE INDEX idx_assets_tenant_created ON assets(tenant_id, created_at DESC);
```

### 4.2 Asset Events

```sql
CREATE TABLE asset_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    asset_id UUID NOT NULL REFERENCES assets(id),
    user_id UUID REFERENCES users(id),
    event_type VARCHAR(50) NOT NULL,
    description TEXT,
    data JSONB,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_asset_events_tenant_id ON asset_events(tenant_id);
CREATE INDEX idx_asset_events_asset_id ON asset_events(asset_id);
CREATE INDEX idx_asset_events_user_id ON asset_events(user_id);
CREATE INDEX idx_asset_events_created_at ON asset_events(created_at DESC);
CREATE INDEX idx_asset_events_data ON asset_events USING GIN(data);
```

### 4.3 Asset Mutations

```sql
CREATE TABLE asset_mutations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    asset_id UUID NOT NULL REFERENCES assets(id),
    from_room_id UUID REFERENCES rooms(id),
    to_room_id UUID REFERENCES rooms(id),
    from_bed_id INTEGER REFERENCES beds(id),
    to_bed_id INTEGER REFERENCES beds(id),
    moved_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_asset_mutations_tenant_id ON asset_mutations(tenant_id);
CREATE INDEX idx_asset_mutations_asset_id ON asset_mutations(asset_id);
CREATE INDEX idx_asset_mutations_from_room_id ON asset_mutations(from_room_id);
CREATE INDEX idx_asset_mutations_to_room_id ON asset_mutations(to_room_id);
CREATE INDEX idx_asset_mutations_created_at ON asset_mutations(created_at DESC);
```

---

## 5. OPERATIONAL TABLES

### 5.1 Complaints

```sql
CREATE TABLE complaints (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    asset_id UUID NOT NULL REFERENCES assets(id),
    reported_by UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'open',
    priority VARCHAR(50) DEFAULT 'normal',
    resolved_at TIMESTAMP,
    resolved_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_complaints_tenant_id ON complaints(tenant_id);
CREATE INDEX idx_complaints_asset_id ON complaints(asset_id);
CREATE INDEX idx_complaints_reported_by ON complaints(reported_by);
CREATE INDEX idx_complaints_status ON complaints(status);
CREATE INDEX idx_complaints_priority ON complaints(priority);
CREATE INDEX idx_complaints_created_at ON complaints(created_at DESC);
CREATE INDEX idx_complaints_deleted_at ON complaints(deleted_at);
```

### 5.2 Maintenance Schedules

```sql
CREATE TABLE maintenance_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    asset_id UUID NOT NULL REFERENCES assets(id),
    schedule_type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    scheduled_date TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'scheduled',
    completed_at TIMESTAMP,
    completed_by UUID REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_maintenance_schedules_tenant_id ON maintenance_schedules(tenant_id);
CREATE INDEX idx_maintenance_schedules_asset_id ON maintenance_schedules(asset_id);
CREATE INDEX idx_maintenance_schedules_scheduled_date ON maintenance_schedules(scheduled_date);
CREATE INDEX idx_maintenance_schedules_status ON maintenance_schedules(status);
CREATE INDEX idx_maintenance_schedules_deleted_at ON maintenance_schedules(deleted_at);
```

### 5.3 Maintenance Documents

```sql
CREATE TABLE maintenance_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    schedule_id UUID NOT NULL REFERENCES maintenance_schedules(id),
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(100),
    uploaded_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_maintenance_documents_tenant_id ON maintenance_documents(tenant_id);
CREATE INDEX idx_maintenance_documents_schedule_id ON maintenance_documents(schedule_id);
CREATE INDEX idx_maintenance_documents_deleted_at ON maintenance_documents(deleted_at);
```

### 5.4 Stock Opname Sessions

```sql
CREATE TABLE stock_opname_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    session_name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'open',
    started_by UUID REFERENCES users(id),
    closed_by UUID REFERENCES users(id),
    started_at TIMESTAMP NOT NULL,
    closed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_stock_opname_sessions_tenant_id ON stock_opname_sessions(tenant_id);
CREATE INDEX idx_stock_opname_sessions_status ON stock_opname_sessions(status);
CREATE INDEX idx_stock_opname_sessions_started_at ON stock_opname_sessions(started_at DESC);
CREATE INDEX idx_stock_opname_sessions_deleted_at ON stock_opname_sessions(deleted_at);
```

### 5.5 Stock Opname Items

```sql
CREATE TABLE stock_opname_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    session_id UUID NOT NULL REFERENCES stock_opname_sessions(id),
    asset_id UUID NOT NULL REFERENCES assets(id),
    condition VARCHAR(50) NOT NULL,
    notes TEXT,
    checked_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_stock_opname_items_tenant_id ON stock_opname_items(tenant_id);
CREATE INDEX idx_stock_opname_items_session_id ON stock_opname_items(session_id);
CREATE INDEX idx_stock_opname_items_asset_id ON stock_opname_items(asset_id);
CREATE INDEX idx_stock_opname_items_condition ON stock_opname_items(condition);
```

---

## 6. DOCUMENT TABLES

### 6.1 Documents

```sql
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    document_type VARCHAR(50) NOT NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_documents_tenant_id ON documents(tenant_id);
CREATE INDEX idx_documents_document_type ON documents(document_type);
CREATE INDEX idx_documents_created_at ON documents(created_at DESC);
CREATE INDEX idx_documents_deleted_at ON documents(deleted_at);
```

### 6.2 Document Files

```sql
CREATE TABLE document_files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(100),
    uploaded_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_document_files_tenant_id ON document_files(tenant_id);
CREATE INDEX idx_document_files_document_id ON document_files(document_id);
CREATE INDEX idx_document_files_deleted_at ON document_files(deleted_at);
```

---

## 7. AUDIT TABLE

### 7.1 Audit Trails

```sql
CREATE TABLE audit_trails (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id),
    user_id UUID REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    table_name VARCHAR(100) NOT NULL,
    record_id VARCHAR(100) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_audit_trails_tenant_id ON audit_trails(tenant_id);
CREATE INDEX idx_audit_trails_user_id ON audit_trails(user_id);
CREATE INDEX idx_audit_trails_action ON audit_trails(action);
CREATE INDEX idx_audit_trails_table_name ON audit_trails(table_name);
CREATE INDEX idx_audit_trails_record_id ON audit_trails(record_id);
CREATE INDEX idx_audit_trails_created_at ON audit_trails(created_at DESC);
CREATE INDEX idx_audit_trails_old_value ON audit_trails USING GIN(old_value);
CREATE INDEX idx_audit_trails_new_value ON audit_trails USING GIN(new_value);
```

---

## 8. RELATIONSHIPS

### 8.1 Entity Relationship Diagram

```
tenants (1) ──< (N) users
tenants (1) ──< (N) roles
tenants (1) ──< (N) assets
tenants (1) ──< (N) rooms
tenants (1) ──< (N) categories

roles (N) ──< (N) permissions (via role_permissions)
users (N) ──> (1) roles

assets (N) ──> (1) categories
assets (N) ──> (1) rooms
assets (N) ──> (1) beds
assets (N) ──> (1) vendors
assets (N) ──> (1) asset_brands
assets (N) ──> (1) asset_models

rooms (1) ──< (N) beds
rooms (1) ──< (N) assets

asset_brands (1) ──< (N) asset_models

assets (1) ──< (N) asset_events
assets (1) ──< (N) asset_mutations
assets (1) ──< (N) complaints
assets (1) ──< (N) maintenance_schedules

maintenance_schedules (1) ──< (N) maintenance_documents

stock_opname_sessions (1) ──< (N) stock_opname_items
assets (1) ──< (N) stock_opname_items

documents (1) ──< (N) document_files
```

---

## 9. INDEXES STRATEGY

### 9.1 Primary Indexes

- Primary keys (automatic)
- Foreign keys
- Unique constraints

### 9.2 Query Optimization Indexes

- `tenant_id` di semua tabel (critical!)
- Composite indexes untuk common queries
- Status fields untuk filtering
- Created_at untuk sorting
- Deleted_at untuk soft delete queries

### 9.3 JSONB Indexes

- GIN indexes untuk JSONB columns (audit_trails, asset_events)

---

## 10. DATA INTEGRITY

### 10.1 Constraints

- NOT NULL untuk required fields
- UNIQUE untuk business keys (code, slug)
- FOREIGN KEY dengan ON DELETE CASCADE/RESTRICT
- CHECK constraints untuk enums

### 10.2 Soft Delete

Semua tabel menggunakan `deleted_at` untuk soft delete:
- Unique indexes dengan `WHERE deleted_at IS NULL`
- Queries filter `deleted_at IS NULL`

---

**Next:** [Authentication & Authorization](./AUTH_RBAC.md)
