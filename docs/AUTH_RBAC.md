# Authentication & RBAC - SAMRS

## 1. AUTHENTICATION

### 1.1 JWT-Based Authentication

SAMRS menggunakan JWT (JSON Web Token) untuk authentication.

**Flow:**
```
1. User login dengan username & password
2. Backend validate credentials
3. Backend generate JWT token
4. Frontend simpan token di localStorage
5. Setiap request, frontend kirim token di header
6. Backend validate token di middleware
7. Extract user info dari token claims
```

### 1.2 JWT Payload

```json
{
  "user_id": "uuid",
  "tenant_id": "uuid",
  "role_id": 1,
  "username": "admin",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### 1.3 Login Flow

**Backend:**
```go
// internal/usecase/auth/auth_usecase.go
func (u *authUsecase) Login(username, password string) (*LoginResponse, error) {
    // 1. Find user by username
    user, err := u.userRepo.GetByUsername(username)
    if err != nil {
        return nil, util.ErrUnauthorized("Invalid credentials")
    }
    
    // 2. Check if user is active
    if !user.IsActive {
        return nil, util.ErrUnauthorized("User is inactive")
    }
    
    // 3. Verify password
    if err := bcrypt.CompareHashAndPassword(
        []byte(user.PasswordHash), 
        []byte(password),
    ); err != nil {
        return nil, util.ErrUnauthorized("Invalid credentials")
    }
    
    // 4. Get user permissions
    permissions, err := u.permissionRepo.ListSlugsByRoleID(user.RoleID)
    if err != nil {
        return nil, err
    }
    
    // 5. Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":   user.ID.String(),
        "tenant_id": user.TenantID.String(),
        "role_id":   user.RoleID,
        "username":  user.Username,
        "exp":       time.Now().Add(24 * time.Hour).Unix(),
        "iat":       time.Now().Unix(),
    })
    
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return nil, err
    }
    
    return &LoginResponse{
        Token:       tokenString,
        User:        user,
        Permissions: permissions,
    }, nil
}
```

**Frontend:**
```typescript
// src/modules/auth/hooks/useLogin.ts
export function useLogin() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [login, { isLoading }] = useLoginMutation();
  
  const handleLogin = async (username: string, password: string) => {
    try {
      const response = await login({ username, password }).unwrap();
      
      // Save token to localStorage
      localStorage.setItem("samrs_token", response.data.token);
      
      // Save user & permissions to Redux
      dispatch(setCredentials({
        user: response.data.user,
        token: response.data.token,
        permissions: response.data.permissions,
      }));
      
      // Redirect to dashboard
      navigate({ to: `/${response.data.user.tenant_slug}` });
      
      toast.success("Login successful");
    } catch (error) {
      toast.error("Invalid credentials");
    }
  };
  
  return { handleLogin, isLoading };
}
```

### 1.4 Auth Middleware

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
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        if err != nil || !token.Valid {
            util.ErrorResponse(c, 401, "Invalid or expired token", nil)
            c.Abort()
            return
        }
        
        // 4. Extract claims & set to context
        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["user_id"])
        c.Set("tenant_id", claims["tenant_id"])
        c.Set("role_id", claims["role_id"])
        c.Set("username", claims["username"])
        
        c.Next()
    }
}
```

---

## 2. RBAC (Role-Based Access Control)

### 2.1 RBAC Model

```
User ──> Role ──> Permissions
```

**Entities:**
- **User**: System user
- **Role**: Group of permissions (e.g., Admin, Technician, Staff)
- **Permission**: Specific action (e.g., `asset:create`, `asset:read`)
- **RolePermission**: Many-to-many relation

### 2.2 Permission Naming Convention

Format: `resource:action`

**Examples:**
- `asset:create` - Create asset
- `asset:read` - Read asset
- `asset:update` - Update asset
- `asset:delete` - Delete asset
- `room:create` - Create room
- `maintenance:complete` - Complete maintenance

### 2.3 Permission List

**Master Data:**
- `room:create`, `room:read`, `room:update`, `room:delete`
- `bed:create`, `bed:read`, `bed:update`, `bed:delete`
- `category:create`, `category:read`, `category:update`, `category:delete`
- `vendor:create`, `vendor:read`, `vendor:update`, `vendor:delete`
- `brand:create`, `brand:read`, `brand:update`, `brand:delete`
- `model:create`, `model:read`, `model:update`, `model:delete`
- `asset_status:create`, `asset_status:read`, `asset_status:update`, `asset_status:delete`

**Asset Management:**
- `asset:create`, `asset:read`, `asset:update`, `asset:delete`
- `asset_mutation:create`, `asset_mutation:read`
- `stock_opname:create`, `stock_opname:read`, `stock_opname:update`, `stock_opname:close`

**Operational:**
- `complaint:create`, `complaint:read`, `complaint:update`, `complaint:delete`
- `maintenance:create`, `maintenance:read`, `maintenance:update`, `maintenance:delete`, `maintenance:complete`
- `maintenance_document:upload`, `maintenance_document:read`, `maintenance_document:delete`

**Document:**
- `document:create`, `document:read`, `document:update`, `document:delete`

**User Management:**
- `user:create`, `user:read`, `user:update`, `user:delete`
- `role:create`, `role:read`, `role:update`, `role:delete`
- `permission:read`, `permission:assign`, `permission:revoke`

**Audit:**
- `audit:read`

**Reporting:**
- `report:export`

**Notification:**
- `notification:send`

**Tenant Management (Super Admin only):**
- `tenant:create`, `tenant:read`, `tenant:update`, `tenant:set_status`

### 2.4 Default Roles

**1. Super Admin (Global)**
- Tenant ID: NULL
- All permissions
- Can manage all tenants
- Can bypass tenant scope

**2. Tenant Admin**
- All permissions within tenant
- Can manage users & roles
- Cannot access other tenants

**3. Technician**
- Asset read/update
- Maintenance read/update/complete
- Complaint read/update
- Stock opname read/update

**4. Staff/Nurse**
- Asset read
- Complaint create/read
- Maintenance read

### 2.5 RBAC Middleware

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

**Usage:**
```go
// internal/app/routes/routes_asset.go
func registerAssetRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
    v1.POST("/assets", 
        requirePerm(deps, "asset:create"),  // RBAC check
        handlers.Asset.Create,
    )
    
    v1.GET("/assets", 
        requirePerm(deps, "asset:read"),
        handlers.Asset.GetAll,
    )
    
    v1.PATCH("/assets/:id", 
        requirePerm(deps, "asset:update"),
        handlers.Asset.Update,
    )
    
    v1.DELETE("/assets/:id", 
        requirePerm(deps, "asset:delete"),
        handlers.Asset.Delete,
    )
}
```

### 2.6 Frontend Permission Check

**Hook:**
```typescript
// src/hooks/usePermission.ts
export function usePermission() {
  const { user } = useAppSelector((s) => s.auth);
  
  const hasPermission = useCallback(
    (permission: string) => {
      if (!user?.permissions) return false;
      return user.permissions.includes(permission);
    },
    [user]
  );
  
  const hasAnyPermission = useCallback(
    (permissions: string[]) => {
      return permissions.some((p) => hasPermission(p));
    },
    [hasPermission]
  );
  
  const hasAllPermissions = useCallback(
    (permissions: string[]) => {
      return permissions.every((p) => hasPermission(p));
    },
    [hasPermission]
  );
  
  return { hasPermission, hasAnyPermission, hasAllPermissions };
}
```

**Component:**
```typescript
// src/modules/auth/components/Can.tsx
export function Can({ 
  permission, 
  children 
}: { 
  permission: string; 
  children: React.ReactNode 
}) {
  const { hasPermission } = usePermission();
  
  if (!hasPermission(permission)) return null;
  
  return <>{children}</>;
}

// Usage
<Can permission={SCOPES.ASSET.CREATE}>
  <Button onClick={handleCreate}>Create Asset</Button>
</Can>
```

**Conditional Rendering:**
```typescript
function AssetListPage() {
  const { hasPermission } = usePermission();
  
  return (
    <PageContainer
      title="Assets"
      actions={
        hasPermission(SCOPES.ASSET.CREATE) && (
          <Button onClick={handleCreate}>
            <Plus className="mr-2 h-4 w-4" />
            Create Asset
          </Button>
        )
      }
    >
      <AssetTable />
    </PageContainer>
  );
}
```

---

## 3. ROLE MANAGEMENT

### 3.1 Create Role

```go
// internal/usecase/rbac/role_usecase.go
func (u *roleUsecase) CreateRole(input CreateRoleInput) (*domain.Role, error) {
    // 1. Validate input
    if input.Name == "" {
        return nil, util.ErrValidation("Role name is required")
    }
    
    // 2. Generate slug
    slug := slug.Make(input.Name)
    
    // 3. Create role
    role := &domain.Role{
        TenantID: input.TenantID,
        Name:     input.Name,
        Slug:     slug,
    }
    
    if err := u.roleRepo.Create(role); err != nil {
        return nil, err
    }
    
    return role, nil
}
```

### 3.2 Assign Permissions

```go
// internal/usecase/rbac/role_permission_usecase.go
func (u *rolePermissionUsecase) AssignPermissions(
    roleID int, 
    permissionIDs []int,
) error {
    // 1. Validate role exists
    _, err := u.roleRepo.FindByID(roleID)
    if err != nil {
        return util.ErrNotFound("Role not found")
    }
    
    // 2. Assign permissions
    if err := u.permissionRepo.AssignPermissions(roleID, permissionIDs); err != nil {
        return err
    }
    
    return nil
}
```

### 3.3 Check Permission

```go
// internal/usecase/rbac/rbac_usecase.go
func (u *rbacUsecase) Authorize(roleID int, permissionSlug string) (bool, error) {
    // 1. Get role
    role, err := u.roleRepo.FindByID(roleID)
    if err != nil {
        return false, err
    }
    
    // 2. Check if role has permission
    hasPermission, err := u.permissionRepo.RoleHasPermission(roleID, permissionSlug)
    if err != nil {
        return false, err
    }
    
    return hasPermission, nil
}
```

---

## 4. TENANT ADMIN

### 4.1 Tenant Admin Middleware

```go
// internal/delivery/http/middleware/tenant_admin_middleware.go
func TenantAdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        roleIDRaw, exists := c.Get("role_id")
        if !exists {
            util.ErrorResponse(c, 401, "Role ID not found", nil)
            c.Abort()
            return
        }
        
        roleID := roleIDRaw.(int)
        
        // Check if role is tenant admin (role_id = 2)
        if roleID != 2 {
            util.ErrorResponse(c, 403, "Tenant admin only", nil)
            c.Abort()
            return
        }
        
        c.Next()
    }
}

func SuperAdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        roleIDRaw, exists := c.Get("role_id")
        if !exists {
            util.ErrorResponse(c, 401, "Role ID not found", nil)
            c.Abort()
            return
        }
        
        roleID := roleIDRaw.(int)
        
        // Check if role is super admin (role_id = 1)
        if roleID != 1 {
            util.ErrorResponse(c, 403, "Super admin only", nil)
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

## 5. SECURITY BEST PRACTICES

### 5.1 Password Security

✅ Hash passwords dengan bcrypt
✅ Minimum password length: 8 characters
✅ Never return password hash to client
✅ Password reset via secure token

### 5.2 JWT Security

✅ Use strong JWT_SECRET (min 32 characters)
✅ Set reasonable expiration (24 hours)
✅ Validate token signature
✅ Check token expiration
✅ Refresh token mechanism (future)

### 5.3 RBAC Security

✅ Always check permission di backend
✅ Frontend permission check hanya untuk UX
✅ Never trust client-side permission check
✅ Audit permission changes

### 5.4 Session Security

✅ Logout clear token dari localStorage
✅ Auto logout on token expiration
✅ Detect concurrent sessions (future)

---

## 6. TESTING

### 6.1 Auth Testing

```go
func TestAuthUsecase_Login(t *testing.T) {
    // Test valid credentials
    // Test invalid credentials
    // Test inactive user
    // Test non-existent user
}
```

### 6.2 RBAC Testing

```go
func TestRBACUsecase_Authorize(t *testing.T) {
    // Test role with permission
    // Test role without permission
    // Test non-existent role
}
```

---

**Next:** [Multi-Tenancy Implementation](./MULTI_TENANCY.md)
