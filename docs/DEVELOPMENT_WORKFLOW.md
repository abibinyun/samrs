# Development Workflow - SAMRS

## 1. SETUP DEVELOPMENT ENVIRONMENT

### 1.1 Prerequisites

**Required:**
- Go 1.25+
- Node.js 18+ (with Bun)
- PostgreSQL 15
- Git
- VS Code (recommended)

**Optional:**
- Docker & Docker Compose
- pgAdmin
- Postman / REST Client extension

### 1.2 Clone Repository

```bash
git clone <repository-url>
cd samrs-project
```

### 1.3 Backend Setup

```bash
cd samrs-backend

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
nano .env

# Start PostgreSQL (via Docker)
docker-compose up -d db

# Install dependencies
go mod download

# Run migrations & seed
go run cmd/seed/main.go

# Start server
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

### 1.4 Frontend Setup

```bash
cd samrs-frontend

# Install dependencies
bun install

# Start BFF & Frontend
bun dev
```

- BFF: `http://localhost:3000`
- Frontend: `http://localhost:5173`

### 1.5 VS Code Extensions

**Recommended:**
- Go (golang.go)
- REST Client (humao.rest-client)
- ESLint (dbaeumer.vscode-eslint)
- Prettier (esbenp.prettier-vscode)
- Tailwind CSS IntelliSense (bradlc.vscode-tailwindcss)
- Error Lens (usernamehw.errorlens)

---

## 2. MENAMBAH FITUR BARU

### 2.1 Backend: Menambah Module Baru

**Contoh: Menambah module "Equipment"**

#### Step 1: Buat Domain Entity

```go
// internal/domain/equipment.go
package domain

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Equipment struct {
    ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    TenantID    uuid.UUID      `gorm:"type:uuid;index;not null"`
    Code        string         `gorm:"size:100;not null"`
    Name        string         `gorm:"size:255;not null"`
    Description string         `gorm:"type:text"`
    Status      string         `gorm:"size:50;default:'active'"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
    
    Tenant      Tenant         `gorm:"foreignKey:TenantID"`
}
```

#### Step 2: Buat Repository

```go
// internal/repository/equipment.go
package repository

import (
    "samrs-backend/internal/domain"
    "github.com/google/uuid"
)

type EquipmentRepository interface {
    Create(equipment *domain.Equipment) error
    FindAllByTenant(tenantID uuid.UUID, filter EquipmentFilter) ([]domain.Equipment, int64, error)
    FindByIDAndTenant(id, tenantID uuid.UUID) (*domain.Equipment, error)
    Update(equipment *domain.Equipment) error
    Delete(id, tenantID uuid.UUID) error
}

type EquipmentFilter struct {
    Page     int
    PageSize int
    Search   string
    Status   string
    SortBy   string
    SortDir  string
}
```

```go
// internal/repository/equipment/equipment_repository.go
package equipment

import (
    "samrs-backend/internal/domain"
    "samrs-backend/internal/repository"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type equipmentRepository struct {
    db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) repository.EquipmentRepository {
    return &equipmentRepository{db: db}
}

func (r *equipmentRepository) Create(equipment *domain.Equipment) error {
    return r.db.Create(equipment).Error
}

func (r *equipmentRepository) FindAllByTenant(
    tenantID uuid.UUID, 
    filter repository.EquipmentFilter,
) ([]domain.Equipment, int64, error) {
    var equipments []domain.Equipment
    var total int64
    
    query := r.db.Model(&domain.Equipment{}).
        Where("tenant_id = ?", tenantID)
    
    // Search
    if filter.Search != "" {
        query = query.Where(
            "code ILIKE ? OR name ILIKE ?",
            "%"+filter.Search+"%",
            "%"+filter.Search+"%",
        )
    }
    
    // Count
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Pagination
    offset := (filter.Page - 1) * filter.PageSize
    query = query.Offset(offset).Limit(filter.PageSize)
    
    // Execute
    if err := query.Find(&equipments).Error; err != nil {
        return nil, 0, err
    }
    
    return equipments, total, nil
}

// ... implement other methods
```

#### Step 3: Buat Usecase

```go
// internal/usecase/equipment/equipment_usecase.go
package usecase

import (
    "samrs-backend/internal/domain"
    "samrs-backend/internal/repository"
    "samrs-backend/pkg/util"
    "github.com/google/uuid"
)

type EquipmentUsecase interface {
    CreateEquipment(input CreateEquipmentInput) (*domain.Equipment, error)
    GetAllEquipments(tenantID uuid.UUID, filter repository.EquipmentFilter) ([]domain.Equipment, int64, error)
    GetEquipmentByID(tenantID uuid.UUID, id uuid.UUID) (*domain.Equipment, error)
    UpdateEquipment(input UpdateEquipmentInput) (*domain.Equipment, error)
    DeleteEquipment(tenantID uuid.UUID, id uuid.UUID) error
}

type CreateEquipmentInput struct {
    TenantID    uuid.UUID
    Code        string
    Name        string
    Description string
    Status      string
}

type equipmentUsecase struct {
    equipmentRepo repository.EquipmentRepository
}

func NewEquipmentUsecase(er repository.EquipmentRepository) EquipmentUsecase {
    return &equipmentUsecase{equipmentRepo: er}
}

func (u *equipmentUsecase) CreateEquipment(input CreateEquipmentInput) (*domain.Equipment, error) {
    // Validate input
    if input.Code == "" {
        return nil, util.ErrValidation("Code is required")
    }
    if input.Name == "" {
        return nil, util.ErrValidation("Name is required")
    }
    
    // Create equipment
    equipment := &domain.Equipment{
        TenantID:    input.TenantID,
        Code:        input.Code,
        Name:        input.Name,
        Description: input.Description,
        Status:      input.Status,
    }
    
    if err := u.equipmentRepo.Create(equipment); err != nil {
        return nil, err
    }
    
    return equipment, nil
}

// ... implement other methods
```

#### Step 4: Buat Handler

```go
// internal/delivery/http/handlers/equipment_handler.go
package handlers

import (
    "net/http"
    "samrs-backend/internal/delivery/http/httputil"
    equipmentusecase "samrs-backend/internal/usecase/equipment"
    "samrs-backend/pkg/util"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type EquipmentHandler struct {
    usecase equipmentusecase.EquipmentUsecase
    db      *gorm.DB
}

func NewEquipmentHandler(u equipmentusecase.EquipmentUsecase, db *gorm.DB) *EquipmentHandler {
    return &EquipmentHandler{usecase: u, db: db}
}

func (h *EquipmentHandler) Create(c *gin.Context) {
    tenantID, ok := httputil.TenantIDFromContext(c)
    if !ok {
        util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID not found", nil)
        return
    }
    
    var input struct {
        Code        string `json:"code" binding:"required"`
        Name        string `json:"name" binding:"required"`
        Description string `json:"description"`
        Status      string `json:"status"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        util.ErrorResponseFromErr(c, "Invalid input", err)
        return
    }
    
    equipment, err := h.usecase.CreateEquipment(equipmentusecase.CreateEquipmentInput{
        TenantID:    tenantID,
        Code:        input.Code,
        Name:        input.Name,
        Description: input.Description,
        Status:      input.Status,
    })
    
    if err != nil {
        util.ErrorResponseFromErr(c, "Failed to create equipment", err)
        return
    }
    
    util.SuccessResponse(c, http.StatusCreated, "Equipment created successfully", equipment)
}

// ... implement other handlers
```

#### Step 5: Register Routes

```go
// internal/app/routes/routes_equipment.go
package routes

import "github.com/gin-gonic/gin"

func registerEquipmentRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
    v1.POST("/equipments", requirePerm(deps, "equipment:create"), handlers.Equipment.Create)
    v1.GET("/equipments", requirePerm(deps, "equipment:read"), handlers.Equipment.GetAll)
    v1.GET("/equipments/:id", requirePerm(deps, "equipment:read"), handlers.Equipment.GetByID)
    v1.PATCH("/equipments/:id", requirePerm(deps, "equipment:update"), handlers.Equipment.Update)
    v1.DELETE("/equipments/:id", requirePerm(deps, "equipment:delete"), handlers.Equipment.Delete)
}
```

#### Step 6: Wire Dependencies

```go
// internal/app/app.go
func New(db *gorm.DB, logger *zap.Logger) *App {
    // ... existing code
    
    // Repository
    equipmentRepo := repository.NewEquipmentRepository(db)
    
    // Usecase
    equipmentUsecase := equipmentusecase.NewEquipmentUsecase(equipmentRepo)
    
    // Handler
    equipmentHandler := handlers.NewEquipmentHandler(equipmentUsecase, db)
    
    // ... existing code
    
    handlers := routes.Handlers{
        // ... existing handlers
        Equipment: equipmentHandler,
    }
    
    // ... existing code
}
```

```go
// internal/app/routes/routes.go
type Handlers struct {
    // ... existing handlers
    Equipment *deliveryhandlers.EquipmentHandler
}
```

```go
// internal/app/routes/routes_protected.go
func registerProtectedRoutes(r *gin.Engine, handlers Handlers, deps RouteDeps) {
    // ... existing routes
    registerEquipmentRoutes(v1, handlers, deps)
}
```

#### Step 7: Add Migration

```go
// internal/config/database.go
err = db.AutoMigrate(
    // ... existing models
    &domain.Equipment{},
)
```

#### Step 8: Add Permissions

```sql
-- Add to seed or migration
INSERT INTO permissions (name, slug, description) VALUES
('Create Equipment', 'equipment:create', 'Can create equipment'),
('Read Equipment', 'equipment:read', 'Can read equipment'),
('Update Equipment', 'equipment:update', 'Can update equipment'),
('Delete Equipment', 'equipment:delete', 'Can delete equipment');
```

---

### 2.2 Frontend: Menambah Module Baru

**Contoh: Menambah module "Equipment"**

#### Step 1: Buat Module Structure

```bash
mkdir -p src/modules/equipment/{actions,components,hooks,layouts,pages}
```

#### Step 2: Define Types

```typescript
// src/modules/equipment/types.ts
export interface Equipment {
  id: string;
  tenant_id: string;
  code: string;
  name: string;
  description?: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface CreateEquipmentInput {
  code: string;
  name: string;
  description?: string;
  status?: string;
}

export interface UpdateEquipmentInput extends CreateEquipmentInput {
  id: string;
}
```

#### Step 3: Define API

```typescript
// src/modules/equipment/actions/equipmentApi.ts
import { api } from "@/store/api";
import type { Equipment, CreateEquipmentInput, UpdateEquipmentInput } from "../types";

export const equipmentApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getEquipments: builder.query<{ data: Equipment[]; meta: any }, void>({
      query: () => "/equipments",
      providesTags: ["Equipment"],
    }),
    
    getEquipmentById: builder.query<{ data: Equipment }, string>({
      query: (id) => `/equipments/${id}`,
      providesTags: (result, error, id) => [{ type: "Equipment", id }],
    }),
    
    createEquipment: builder.mutation<{ data: Equipment }, CreateEquipmentInput>({
      query: (input) => ({
        url: "/equipments",
        method: "POST",
        body: input,
      }),
      invalidatesTags: ["Equipment"],
    }),
    
    updateEquipment: builder.mutation<{ data: Equipment }, UpdateEquipmentInput>({
      query: ({ id, ...input }) => ({
        url: `/equipments/${id}`,
        method: "PATCH",
        body: input,
      }),
      invalidatesTags: (result, error, { id }) => [{ type: "Equipment", id }],
    }),
    
    deleteEquipment: builder.mutation<void, string>({
      query: (id) => ({
        url: `/equipments/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["Equipment"],
    }),
  }),
});

export const {
  useGetEquipmentsQuery,
  useGetEquipmentByIdQuery,
  useCreateEquipmentMutation,
  useUpdateEquipmentMutation,
  useDeleteEquipmentMutation,
} = equipmentApi;
```

#### Step 4: Create Pages

```typescript
// src/modules/equipment/pages/EquipmentListPage.tsx
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { useNavigate } from "@tanstack/react-router";
import { useGetEquipmentsQuery } from "../actions/equipmentApi";
import { usePermission } from "@/hooks/usePermission";
import { SCOPES } from "@/constants/permissions";

export default function EquipmentListPage() {
  const navigate = useNavigate();
  const { hasPermission } = usePermission();
  const { data, isLoading } = useGetEquipmentsQuery();
  
  const equipments = data?.data ?? [];
  
  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Equipments</h1>
        {hasPermission(SCOPES.EQUIPMENT.CREATE) && (
          <Button onClick={() => navigate({ to: "/equipments/create" })}>
            <Plus className="mr-2 h-4 w-4" />
            Create Equipment
          </Button>
        )}
      </div>
      
      {isLoading ? (
        <div>Loading...</div>
      ) : (
        <div>
          {/* Table or list of equipments */}
        </div>
      )}
    </div>
  );
}
```

#### Step 5: Create Routes

```typescript
// src/modules/equipment/routes.tsx
import { createRoute, lazyRouteComponent } from "@tanstack/react-router";
import { tenantRoute } from "@/routes";

export const equipmentRoutes = createRoute({
  getParentRoute: () => tenantRoute,
  path: "/equipments",
  component: lazyRouteComponent(() => import("./pages/EquipmentListPage")),
});

export const equipmentCreateRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "/equipments/create",
  component: lazyRouteComponent(() => import("./pages/EquipmentCreatePage")),
});
```

#### Step 6: Register Routes

```typescript
// src/routes/index.tsx
import { equipmentRoutes, equipmentCreateRoute } from "@/modules/equipment/routes";

const routeTree = rootRoute.addChildren([
  loginRoute,
  tenantRoute.addChildren([
    dashboardRoute,
    assetRoutes,
    equipmentRoutes,
    equipmentCreateRoute,
    // ... other routes
  ]),
]);
```

#### Step 7: Add Permissions

```typescript
// src/constants/permissions.ts
export const SCOPES = {
  // ... existing scopes
  EQUIPMENT: {
    CREATE: "equipment:create",
    READ: "equipment:read",
    UPDATE: "equipment:update",
    DELETE: "equipment:delete",
  },
};
```

#### Step 8: Add Navigation

```typescript
// src/constants/navigations.ts
export const MENU_ITEMS: MenuSection[] = [
  // ... existing sections
  {
    group: "Equipment Management",
    items: [
      {
        label: "Equipments",
        icon: Wrench,
        to: "/equipments",
        permission: SCOPES.EQUIPMENT.READ,
      },
    ],
  },
];
```

---

## 3. TESTING

### 3.1 Backend Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/usecase/equipment/...

# Run with verbose
go test -v ./...
```

### 3.2 Generate Mocks

```bash
# Install mockgen
go install go.uber.org/mock/mockgen@latest

# Generate mocks
mockgen -source=internal/repository/equipment.go \
        -destination=internal/test/mocks/equipment_repository_mock.go \
        -package=mocks
```

### 3.3 Write Unit Tests

```go
// internal/usecase/equipment/equipment_usecase_test.go
package usecase

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "go.uber.org/mock/gomock"
    "samrs-backend/internal/test/mocks"
)

func TestEquipmentUsecase_CreateEquipment(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockEquipmentRepository(ctrl)
    usecase := NewEquipmentUsecase(mockRepo)
    
    t.Run("success", func(t *testing.T) {
        input := CreateEquipmentInput{
            Code: "EQ-001",
            Name: "Test Equipment",
        }
        
        mockRepo.EXPECT().
            Create(gomock.Any()).
            Return(nil)
        
        equipment, err := usecase.CreateEquipment(input)
        
        assert.NoError(t, err)
        assert.NotNil(t, equipment)
        assert.Equal(t, input.Code, equipment.Code)
    })
}
```

---

## 4. GIT WORKFLOW

### 4.1 Branch Strategy

```
main (production)
  ↓
develop (staging)
  ↓
feature/nama-fitur (development)
```

### 4.2 Create Feature Branch

```bash
# Update develop
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feature/equipment-module

# Work on feature
# ... make changes

# Commit changes
git add .
git commit -m "feat: add equipment module"

# Push to remote
git push origin feature/equipment-module

# Create Pull Request to develop
```

### 4.3 Commit Convention

```
feat: new feature
fix: bug fix
docs: documentation
refactor: code refactoring
test: add tests
chore: maintenance
style: formatting
perf: performance improvement
```

**Examples:**
```bash
git commit -m "feat: add equipment CRUD endpoints"
git commit -m "fix: resolve tenant scope issue in asset query"
git commit -m "docs: update API documentation"
git commit -m "refactor: simplify auth middleware"
git commit -m "test: add unit tests for equipment usecase"
```

---

## 5. CODE REVIEW CHECKLIST

### 5.1 Backend

- [ ] Semua query tenant-scoped
- [ ] Permission check di setiap endpoint
- [ ] Input validation
- [ ] Error handling
- [ ] Transaction untuk operasi multi-table
- [ ] Unit tests
- [ ] No hardcoded values
- [ ] Proper logging
- [ ] Documentation

### 5.2 Frontend

- [ ] TypeScript types defined
- [ ] Permission check di UI
- [ ] Loading & error states
- [ ] Form validation
- [ ] Responsive design
- [ ] Accessible components
- [ ] No console.log
- [ ] Proper error handling

---

## 6. DEBUGGING

### 6.1 Backend Debugging

**VS Code launch.json:**
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch API",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/api/main.go",
      "env": {},
      "args": []
    }
  ]
}
```

**Logging:**
```go
// Use structured logging
logger.Info("Creating asset",
    zap.String("code", input.Code),
    zap.String("tenant_id", tenantID.String()),
)
```

### 6.2 Frontend Debugging

**Browser DevTools:**
- Network tab untuk API calls
- Redux DevTools untuk state
- React DevTools untuk components

**Console Logging:**
```typescript
console.log("Asset data:", asset);
console.table(assets);
```

---

**Next:** [API Documentation](./API_DOCUMENTATION.md)
