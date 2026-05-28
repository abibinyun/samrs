# Backend Deep Dive - SAMRS

## 1. ARSITEKTUR BACKEND

### 1.1 Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                    DELIVERY LAYER                        │
│  - HTTP Handlers (Controllers)                          │
│  - Middleware (Auth, RBAC, Logging, etc)                │
│  - Request/Response mapping                             │
│  - Input validation                                     │
└─────────────────────────────────────────────────────────┘
                         ↓ depends on
┌─────────────────────────────────────────────────────────┐
│                    USECASE LAYER                         │
│  - Business Logic                                       │
│  - Orchestration                                        │
│  - Business Rules Validation                            │
│  - Transaction Management                               │
└─────────────────────────────────────────────────────────┘
                         ↓ depends on
┌─────────────────────────────────────────────────────────┐
│                   REPOSITORY LAYER                       │
│  - Data Access Interface                                │
│  - GORM Queries                                         │
│  - Tenant Scoping                                       │
│  - Pagination & Filtering                               │
└─────────────────────────────────────────────────────────┘
                         ↓ depends on
┌─────────────────────────────────────────────────────────┐
│                    DOMAIN LAYER                          │
│  - Entities (Structs)                                   │
│  - Business Constants                                   │
│  - Domain Validation Functions                          │
│  - NO DEPENDENCIES                                      │
└─────────────────────────────────────────────────────────┘
```

### 1.2 Dependency Injection

Semua dependencies di-wire di `internal/app/app.go`:

```go
func New(db *gorm.DB, logger *zap.Logger) *App {
    // 1. Initialize Repositories
    assetRepo := repository.NewAssetRepository(db)
    userRepo := repository.NewUserRepository(db)
    // ... other repos
    
    // 2. Initialize Usecases (inject repos)
    assetUsecase := usecase.NewAssetUsecase(assetRepo, categoryRepo, ...)
    authUsecase := usecase.NewAuthUsecase(userRepo, permissionRepo)
    // ... other usecases
    
    // 3. Initialize Handlers (inject usecases)
    assetHandler := handlers.NewAssetHandler(assetUsecase, db)
    authHandler := handlers.NewAuthHandler(authUsecase)
    // ... other handlers
    
    // 4. Setup Router & Middleware
    router := gin.New()
    router.Use(middleware.RequestID(), middleware.Logger(), ...)
    
    // 5. Register Routes
    routes.Register(router, handlers, deps)
    
    return &App{Router: router}
}
```

---

## 2. DOMAIN LAYER

### 2.1 Domain Entities

Total: **20 domain models**

**Core Identity & RBAC:**
- `Tenant` - Multi-tenant isolation
- `User` - System users
- `Role` - User roles
- `Permission` - Access permissions
- `RolePermission` - Many-to-many relation

**Master Data:**
- `Room` - Lokasi ruangan
- `Bed` - Tempat tidur di ruangan
- `Category` - Kategori aset
- `Vendor` - Vendor/supplier
- `AssetBrand` - Brand aset
- `AssetModel` - Model aset
- `AssetStatus` - Status aset (custom)

**Asset Management:**
- `Asset` - Aset utama
- `AssetEvent` - Timeline events
- `AssetMutation` - Perpindahan aset

**Operational:**
- `Complaint` - Keluhan/laporan
- `MaintenanceSchedule` - Jadwal maintenance
- `MaintenanceDocument` - Dokumen maintenance
- `StockOpnameSession` - Sesi stock opname
- `StockOpnameItem` - Item stock opname
- `Document` - Dokumen umum
- `DocumentFile` - File dokumen

**Audit:**
- `AuditTrail` - Audit log

### 2.2 Contoh Domain Entity

```go
// internal/domain/asset.go
type Asset struct {
    ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    TenantID     uuid.UUID      `gorm:"type:uuid;index;not null"`
    CategoryID   uint           `gorm:"not null"`
    RoomID       *uuid.UUID     `gorm:"type:uuid"`
    BedID        *uint
    VendorID     *uint
    BrandID      *uint
    ModelID      *uint
    Code         string         `gorm:"size:100;not null;index"`
    Name         string         `gorm:"size:255;not null"`
    Status       string         `gorm:"size:50;default:'ready'"`
    PurchaseDate *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
    
    // Relations
    Category Category   `gorm:"foreignKey:CategoryID"`
    Room     Room       `gorm:"foreignKey:RoomID"`
    Bed      Bed        `gorm:"foreignKey:BedID"`
    Vendor   Vendor     `gorm:"foreignKey:VendorID"`
    BrandRef AssetBrand `gorm:"foreignKey:BrandID"`
    ModelRef AssetModel `gorm:"foreignKey:ModelID"`
    Tenant   Tenant     `gorm:"foreignKey:TenantID"`
}

// Domain validation
func IsValidAssetStatus(status string) bool {
    switch status {
    case AssetStatusReady, AssetStatusBroken, AssetStatusMaintenance:
        return true
    default:
        return false
    }
}
```

**Key Points:**
- Semua entity punya `TenantID` (kecuali Tenant & Permission)
- Menggunakan UUID untuk primary key (security)
- Soft delete dengan `DeletedAt`
- GORM tags untuk database mapping
- Relations untuk eager loading

---

## 3. REPOSITORY LAYER

### 3.1 Repository Pattern

Setiap domain entity punya repository interface & implementation.

**Interface Example:**
```go
// internal/repository/asset.go
type AssetRepository interface {
    Create(asset *domain.Asset) error
    FindAllByTenant(tenantID uuid.UUID, filter AssetFilter) ([]domain.Asset, int64, error)
    FindByIDAndTenant(id, tenantID uuid.UUID) (*domain.Asset, error)
    Update(asset *domain.Asset) error
    Delete(id, tenantID uuid.UUID) error
    ExistsByCode(code string, tenantID uuid.UUID) (bool, error)
    CountByTenant(tenantID uuid.UUID) (int64, error)
}

type AssetFilter struct {
    Page     int
    PageSize int
    Search   string
    Status   string
    Category uint
    SortBy   string
    SortDir  string
}
```

**Implementation Example:**
```go
// internal/repository/asset/asset_repository.go
type assetRepository struct {
    db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) repository.AssetRepository {
    return &assetRepository{db: db}
}

func (r *assetRepository) FindAllByTenant(
    tenantID uuid.UUID, 
    filter repository.AssetFilter,
) ([]domain.Asset, int64, error) {
    var assets []domain.Asset
    var total int64
    
    query := r.db.Model(&domain.Asset{}).
        Where("tenant_id = ?", tenantID)
    
    // Search
    if filter.Search != "" {
        query = query.Where(
            "code ILIKE ? OR name ILIKE ?",
            "%"+filter.Search+"%",
            "%"+filter.Search+"%",
        )
    }
    
    // Filter by status
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    
    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Sorting
    sortBy := mapAssetSortBy(filter.SortBy)
    sortDir := base.MapSortDir(filter.SortDir)
    query = query.Order(sortBy + " " + sortDir)
    
    // Pagination
    offset := (filter.Page - 1) * filter.PageSize
    query = query.Offset(offset).Limit(filter.PageSize)
    
    // Preload relations
    query = query.Preload("Category").
        Preload("Room").
        Preload("Bed").
        Preload("Vendor")
    
    if err := query.Find(&assets).Error; err != nil {
        return nil, 0, err
    }
    
    return assets, total, nil
}
```

**Key Points:**
- Repository hanya handle data access
- Tenant scoping di setiap query
- Pagination & filtering
- Preload relations untuk menghindari N+1 query
- Error handling

### 3.2 Tenant Scope Guard

GORM plugin untuk memastikan semua query tenant-scoped:

```go
// internal/config/tenant_scope_guard.go
func tenantScopeGuard(db *gorm.DB) {
    // Check if query has tenant condition
    if !hasTenantCondition(db) {
        // Check if explicitly skipped
        if _, ok := db.Statement.Settings.Load("skip_tenant_scope"); ok {
            return
        }
        
        // PANIC if no tenant scope!
        panic("Query without tenant scope detected!")
    }
}
```

**Cara bypass (hanya untuk super admin):**
```go
db.Session(&gorm.Session{
    Context: config.WithSkipTenantScope(ctx),
})
```

---

## 4. USECASE LAYER

### 4.1 Usecase Pattern

Business logic & orchestration.

**Interface Example:**
```go
// internal/usecase/asset/asset_usecase.go
type AssetUsecase interface {
    CreateAsset(input CreateAssetInput) (*domain.Asset, error)
    GetAllAssets(tenantID uuid.UUID, filter repository.AssetFilter) ([]domain.Asset, int64, error)
    GetAssetByID(tenantID uuid.UUID, id uuid.UUID) (*domain.Asset, error)
    UpdateAsset(input UpdateAssetInput) (*domain.Asset, error)
    DeleteAsset(tenantID uuid.UUID, id uuid.UUID) error
}

type CreateAssetInput struct {
    TenantID     uuid.UUID
    CategoryID   uint
    RoomID       *uuid.UUID
    BedID        *uint
    Code         string
    Name         string
    Status       string
    PurchaseDate string
}
```

**Implementation Example:**
```go
type assetUsecase struct {
    assetRepo    repository.AssetRepository
    categoryRepo repository.CategoryRepository
    roomRepo     repository.RoomRepository
    bedRepo      repository.BedRepository
}

func (u *assetUsecase) CreateAsset(input CreateAssetInput) (*domain.Asset, error) {
    // 1. Validate business rules
    if !domain.IsValidAssetStatus(input.Status) {
        return nil, util.ErrValidation("Invalid asset status")
    }
    
    // 2. Check if code already exists
    exists, err := u.assetRepo.ExistsByCode(input.Code, input.TenantID)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, util.ErrConflict("Asset code already exists")
    }
    
    // 3. Validate category exists
    category, err := u.categoryRepo.FindByID(input.CategoryID)
    if err != nil {
        return nil, util.ErrNotFound("Category not found")
    }
    
    // 4. Validate room if provided
    if input.RoomID != nil {
        _, err := u.roomRepo.FindByIDAndTenant(*input.RoomID, input.TenantID)
        if err != nil {
            return nil, util.ErrNotFound("Room not found")
        }
    }
    
    // 5. Parse purchase date
    var purchaseDate *time.Time
    if input.PurchaseDate != "" {
        parsed, err := time.Parse("2006-01-02", input.PurchaseDate)
        if err != nil {
            return nil, util.ErrValidation("Invalid purchase date format")
        }
        purchaseDate = &parsed
    }
    
    // 6. Create asset
    asset := &domain.Asset{
        TenantID:     input.TenantID,
        CategoryID:   input.CategoryID,
        RoomID:       input.RoomID,
        BedID:        input.BedID,
        Code:         input.Code,
        Name:         input.Name,
        Status:       input.Status,
        PurchaseDate: purchaseDate,
    }
    
    if err := u.assetRepo.Create(asset); err != nil {
        return nil, err
    }
    
    return asset, nil
}
```

**Key Points:**
- Validate business rules
- Orchestrate multiple repository calls
- Handle transactions if needed
- Return custom errors
- No HTTP concerns (pure business logic)

### 4.2 Transaction Management

Untuk operasi multi-table, gunakan transaction:

```go
func (u *assetUsecase) MoveAsset(input MoveAssetInput) error {
    return u.db.Transaction(func(tx *gorm.DB) error {
        // 1. Update asset location
        asset, err := u.assetRepo.FindByID(input.AssetID)
        if err != nil {
            return err
        }
        
        asset.RoomID = input.NewRoomID
        asset.BedID = input.NewBedID
        
        if err := u.assetRepo.Update(asset); err != nil {
            return err
        }
        
        // 2. Create mutation record
        mutation := &domain.AssetMutation{
            AssetID:    input.AssetID,
            FromRoomID: input.OldRoomID,
            ToRoomID:   input.NewRoomID,
            MovedBy:    input.UserID,
        }
        
        if err := u.mutationRepo.Create(mutation); err != nil {
            return err
        }
        
        // 3. Create event
        event := &domain.AssetEvent{
            AssetID: input.AssetID,
            Type:    "MOVED",
            Data:    map[string]interface{}{...},
        }
        
        if err := u.eventRepo.Create(event); err != nil {
            return err
        }
        
        return nil
    })
}
```

---

## 5. DELIVERY LAYER

### 5.1 HTTP Handlers

**Handler Structure:**
```go
// internal/delivery/http/handlers/asset_handler.go
type AssetHandler struct {
    usecase assetusecase.AssetUsecase
    db      *gorm.DB
}

func NewAssetHandler(u assetusecase.AssetUsecase, db *gorm.DB) *AssetHandler {
    return &AssetHandler{usecase: u, db: db}
}

func (h *AssetHandler) Create(c *gin.Context) {
    // 1. Get tenant & user from context
    tenantID, ok := httputil.TenantIDFromContext(c)
    if !ok {
        util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID not found", nil)
        return
    }
    
    // 2. Parse & validate input
    var input struct {
        CategoryID   uint   `json:"category_id" binding:"required"`
        Code         string `json:"code" binding:"required"`
        Name         string `json:"name" binding:"required"`
        Status       string `json:"status"`
        PurchaseDate string `json:"purchase_date"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        util.ErrorResponseFromErr(c, "Invalid input", err)
        return
    }
    
    // 3. Call usecase
    asset, err := h.usecase.CreateAsset(assetusecase.CreateAssetInput{
        TenantID:     tenantID,
        CategoryID:   input.CategoryID,
        Code:         input.Code,
        Name:         input.Name,
        Status:       input.Status,
        PurchaseDate: input.PurchaseDate,
    })
    
    if err != nil {
        util.ErrorResponseFromErr(c, "Failed to create asset", err)
        return
    }
    
    // 4. Return response
    util.SuccessResponse(c, http.StatusCreated, "Asset created successfully", asset)
}
```

**Key Points:**
- Extract context (tenant, user)
- Parse & validate input
- Call usecase
- Handle errors
- Return standardized response

### 5.2 Middleware Chain

```go
// internal/app/routes/routes_protected.go
func registerProtectedRoutes(r *gin.Engine, handlers Handlers, deps RouteDeps) {
    v1 := r.Group("/api/v1")
    
    // Apply middleware
    v1.Use(
        middleware.AuthMiddleware(),           // JWT validation
        middleware.TenantScope(),              // Inject tenant to context
        middleware.RateLimiter(limit, burst),  // Rate limiting
    )
    
    // Register routes with RBAC
    v1.POST("/assets", 
        requirePerm(deps, "asset:create"),  // RBAC check
        handlers.Asset.Create,
    )
}
```

**Middleware Order:**
1. Request ID
2. Structured Logger
3. Recovery
4. Metrics
5. Rate Limiter
6. Auth (JWT)
7. Tenant Scope
8. RBAC (per route)

---

## 6. MIDDLEWARE DETAIL

### 6.1 Auth Middleware

```go
// internal/delivery/http/middleware/auth_middleware.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extract token from header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            util.ErrorResponse(c, 401, "Authorization required", nil)
            c.Abort()
            return
        }
        
        // 2. Parse Bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            util.ErrorResponse(c, 401, "Invalid authorization format", nil)
            c.Abort()
            return
        }
        
        // 3. Validate JWT
        token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        if err != nil || !token.Valid {
            util.ErrorResponse(c, 401, "Invalid or expired token", nil)
            c.Abort()
            return
        }
        
        // 4. Extract claims
        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["user_id"])
        c.Set("tenant_id", claims["tenant_id"])
        c.Set("role_id", claims["role_id"])
        
        c.Next()
    }
}
```

### 6.2 RBAC Middleware

```go
// internal/delivery/http/middleware/rbac_middleware.go
func RBACMiddleware(rbac rbacusecase.RBACUsecase, permissionSlug string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Get role from context
        roleIDRaw, exists := c.Get("role_id")
        if !exists {
            util.ErrorResponse(c, 401, "Role ID not found", nil)
            c.Abort()
            return
        }
        
        roleID := roleIDRaw.(int)
        
        // 2. Check permission
        allowed, err := rbac.Authorize(roleID, permissionSlug)
        if err != nil {
            util.ErrorResponse(c, 401, "Role not found", nil)
            c.Abort()
            return
        }
        
        if !allowed {
            util.ErrorResponse(c, 403, "Access denied", nil)
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### 6.3 Tenant Scope Middleware

```go
// internal/delivery/http/middleware/tenant_scope_middleware.go
func TenantScope() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantIDRaw, exists := c.Get("tenant_id")
        if !exists {
            c.Next()
            return
        }
        
        tenantID := tenantIDRaw.(string)
        
        // Inject tenant to GORM context
        ctx := config.MarkTenantScope(c.Request.Context(), tenantID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}
```

---

## 7. ERROR HANDLING

### 7.1 Custom Error Types

```go
// pkg/util/app_error.go
type AppError struct {
    Kind    string
    Message string
    Err     error
}

func ErrValidation(msg string) *AppError {
    return &AppError{Kind: "validation", Message: msg}
}

func ErrNotFound(msg string) *AppError {
    return &AppError{Kind: "not_found", Message: msg}
}

func ErrConflict(msg string) *AppError {
    return &AppError{Kind: "conflict", Message: msg}
}

func ErrUnauthorized(msg string) *AppError {
    return &AppError{Kind: "unauthorized", Message: msg}
}

func ErrForbidden(msg string) *AppError {
    return &AppError{Kind: "forbidden", Message: msg}
}
```

### 7.2 Error Response Mapping

```go
// pkg/util/error_mapping.go
func HTTPStatusFromError(err error) int {
    appErr, ok := err.(*AppError)
    if !ok {
        return http.StatusInternalServerError
    }
    
    switch appErr.Kind {
    case "validation":
        return http.StatusBadRequest
    case "not_found":
        return http.StatusNotFound
    case "conflict":
        return http.StatusConflict
    case "unauthorized":
        return http.StatusUnauthorized
    case "forbidden":
        return http.StatusForbidden
    default:
        return http.StatusInternalServerError
    }
}

func ErrorResponseFromErr(c *gin.Context, message string, err error) {
    status := HTTPStatusFromError(err)
    details := errorDetailsForResponse(err)
    
    c.JSON(status, gin.H{
        "success": false,
        "message": message,
        "error":   details,
    })
}
```

---

## 8. RESPONSE FORMAT

### 8.1 Success Response

```go
// pkg/util/response.go
func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
    c.JSON(status, gin.H{
        "success": true,
        "message": message,
        "data":    data,
    })
}

func SuccessResponseWithMeta(c *gin.Context, status int, message string, data interface{}, meta interface{}) {
    c.JSON(status, gin.H{
        "success": true,
        "message": message,
        "data":    data,
        "meta":    meta,
    })
}
```

**Example:**
```json
{
  "success": true,
  "message": "Assets retrieved successfully",
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

### 8.2 Error Response

```json
{
  "success": false,
  "message": "Failed to create asset",
  "error": {
    "type": "validation",
    "details": "Asset code already exists"
  }
}
```

---

## 9. TESTING

### 9.1 Unit Test Example

```go
// internal/usecase/asset/asset_usecase_test.go
func TestAssetUsecase_CreateAsset(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    // Setup mocks
    mockAssetRepo := mocks.NewMockAssetRepository(ctrl)
    mockCategoryRepo := mocks.NewMockCategoryRepository(ctrl)
    
    usecase := NewAssetUsecase(mockAssetRepo, mockCategoryRepo, ...)
    
    // Test case
    t.Run("success", func(t *testing.T) {
        input := CreateAssetInput{
            TenantID:   uuid.New(),
            CategoryID: 1,
            Code:       "AST-001",
            Name:       "Test Asset",
            Status:     "ready",
        }
        
        // Mock expectations
        mockAssetRepo.EXPECT().
            ExistsByCode(input.Code, input.TenantID).
            Return(false, nil)
        
        mockCategoryRepo.EXPECT().
            FindByID(input.CategoryID).
            Return(&domain.Category{ID: 1}, nil)
        
        mockAssetRepo.EXPECT().
            Create(gomock.Any()).
            Return(nil)
        
        // Execute
        asset, err := usecase.CreateAsset(input)
        
        // Assert
        assert.NoError(t, err)
        assert.NotNil(t, asset)
        assert.Equal(t, input.Code, asset.Code)
    })
}
```

---

## 10. BEST PRACTICES

### 10.1 DO's

✅ Selalu tenant-scope semua query
✅ Gunakan transaction untuk operasi multi-table
✅ Validate input di handler, business rules di usecase
✅ Return custom error types
✅ Preload relations untuk menghindari N+1
✅ Gunakan UUID untuk primary key
✅ Soft delete dengan DeletedAt
✅ Index kolom yang sering di-query
✅ Pagination untuk list endpoint
✅ Structured logging dengan context

### 10.2 DON'Ts

❌ Jangan query tanpa tenant scope
❌ Jangan skip RBAC check
❌ Jangan expose internal error ke client
❌ Jangan hardcode credentials
❌ Jangan bypass audit trail
❌ Jangan lupa handle transaction rollback
❌ Jangan return password hash ke client
❌ Jangan trust user input tanpa validasi

---

**Next:** [Frontend Deep Dive](./FRONTEND_DEEP_DIVE.md)
