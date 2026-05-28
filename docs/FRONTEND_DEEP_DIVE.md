# Frontend Deep Dive - SAMRS

## 1. ARSITEKTUR FRONTEND

### 1.1 Tech Stack Overview

```
React 19 (UI Framework)
  ↓
TypeScript (Type Safety)
  ↓
Vite/Rolldown (Build Tool)
  ↓
Redux Toolkit (State Management)
  ↓
RTK Query (API Client & Caching)
  ↓
TanStack Router (Type-safe Routing)
  ↓
TanStack Table (Headless Table)
  ↓
Shadcn/UI + Tailwind (UI Components)
  ↓
React Hook Form + Zod (Forms & Validation)
```

### 1.2 Application Flow

```
User Action
  ↓
Component Event Handler
  ↓
Dispatch Redux Action / Call RTK Query Hook
  ↓
BFF (Hono Proxy)
  ↓
Backend API
  ↓
Response
  ↓
RTK Query Cache Update
  ↓
Component Re-render
```

---

## 2. STRUKTUR MODULAR

### 2.1 Module Structure

Setiap feature module mengikuti struktur ini:

```
modules/[module-name]/
├── actions/              # Redux & API
│   ├── [module]Api.ts    # RTK Query endpoints
│   ├── [module]Slice.ts  # Redux slice (optional)
│   └── [module]UiSlice.ts # UI state (optional)
├── components/           # Module-specific components
│   └── [Component].tsx
├── hooks/                # Custom hooks
│   └── use[Hook].ts
├── layouts/              # Layout components
│   └── [Layout].tsx
├── pages/                # Page components
│   └── [Page].tsx
├── mock/                 # Mock data (development)
│   └── mocks.ts
├── routes.tsx            # Module routes
├── types.ts              # TypeScript types
└── schemas.ts            # Zod schemas
```

### 2.2 Contoh Module: Asset

```
modules/asset/
├── actions/
│   ├── assetApi.ts       # API endpoints
│   ├── assetSlice.ts     # Global asset state
│   └── assetListUiSlice.ts # UI state per tab
├── components/
│   ├── AssetColumns.tsx  # Table columns
│   ├── AssetStatus.tsx   # Status badge
│   └── AssetQrDialog.tsx # QR dialog
├── hooks/
│   └── useAssetCreate.ts # Create asset logic
├── layouts/
│   ├── AssetTable.tsx    # Table layout
│   └── AssetFormFields.tsx # Form fields
├── pages/
│   ├── AssetListPage.tsx # List page
│   └── AssetCreatePage.tsx # Create page
├── mock/
│   ├── assets.mock.ts    # Mock data
│   ├── mocks.ts          # Worker runner
│   └── worker.ts         # Web worker
├── routes.tsx            # Asset routes
├── types.ts              # Asset types
└── schemas.ts            # Validation schemas
```

---

## 3. STATE MANAGEMENT

### 3.1 Redux Store Structure

```typescript
// src/store/index.ts
export const store = configureStore({
  reducer: {
    // Auth state
    auth: authReducer,
    
    // RTK Query APIs
    [api.reducerPath]: api.reducer,
    [assetApi.reducerPath]: assetApi.reducer,
    [mutationApi.reducerPath]: mutationApi.reducer,
    [formTemplatesApi.reducerPath]: formTemplatesApi.reducer,
    
    // UI state
    assets: assetReducer,
    tabs: tabsReducer,
    tableState: tableStateReducer,
    assetListUi: assetListUiReducer,
    mutations: mutationReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }).concat(
      api.middleware,
      assetApi.middleware,
      mutationApi.middleware,
      formTemplatesApi.middleware
    ),
});
```

### 3.2 RTK Query API

**Base API:**
```typescript
// src/store/api.ts
const rawBaseQuery = fetchBaseQuery({
  baseUrl: import.meta.env.VITE_API_URL,
  credentials: "include",
  prepareHeaders: (headers, { getState }) => {
    const state = getState() as RootState;
    
    // Add JWT token
    if (state.auth.token) {
      headers.set("Authorization", `Bearer ${state.auth.token}`);
    }
    
    // Add tenant header
    if (state.auth.user?.tenant_slug) {
      headers.set("X-Tenant-ID", state.auth.user.tenant_slug);
    }
    
    return headers;
  },
});

export const api = createApi({
  reducerPath: "api",
  baseQuery: rawBaseQuery,
  tagTypes: ["User", "Asset", "Room", ...],
  endpoints: () => ({}),
});
```

**Module API:**
```typescript
// src/modules/asset/actions/assetApi.ts
export const assetApi = createApi({
  reducerPath: "assetApi",
  baseQuery: async () => ({ data: {} }),
  tagTypes: ["Assets"],
  endpoints: (builder) => ({
    getAssets: builder.query<AssetDevice[], void>({
      async queryFn() {
        // Mock implementation
        const data = await runAssetWorker();
        return { data };
      },
      providesTags: ["Assets"],
    }),
    
    createAsset: builder.mutation<Asset, CreateAssetInput>({
      query: (input) => ({
        url: "/assets",
        method: "POST",
        body: input,
      }),
      invalidatesTags: ["Assets"],
    }),
  }),
});

export const { useGetAssetsQuery, useCreateAssetMutation } = assetApi;
```

**Usage in Component:**
```typescript
function AssetListPage() {
  const { data: assets = [], isLoading } = useGetAssetsQuery();
  const [createAsset, { isLoading: isCreating }] = useCreateAssetMutation();
  
  const handleCreate = async (input: CreateAssetInput) => {
    try {
      await createAsset(input).unwrap();
      toast.success("Asset created");
    } catch (error) {
      toast.error("Failed to create asset");
    }
  };
  
  return (
    <div>
      {isLoading ? <Skeleton /> : <AssetTable data={assets} />}
    </div>
  );
}
```

---

## 4. ROUTING

### 4.1 Route Structure

```typescript
// src/routes/index.tsx
import { createRouter } from "@tanstack/react-router";
import { rootRoute } from "./__root";

// Public routes
const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  component: lazyRouteComponent(() => import("../modules/auth/pages/LoginPage")),
});

// Protected tenant route
export const tenantRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/$tenantId",
  beforeLoad: ({ params }) => {
    const { user } = store.getState().auth;
    const token = localStorage.getItem("samrs_token");
    
    if (!token) throw redirect({ to: "/login" });
    if (user?.tenant_slug !== params.tenantId) {
      throw redirect({ to: `/${user?.tenant_slug}` });
    }
  },
  component: DashboardShell,
});

// Module routes
const assetRoutes = createRoute({
  getParentRoute: () => tenantRoute,
  path: "/assets",
  component: AssetListPage,
});

// Build router tree
const routeTree = rootRoute.addChildren([
  loginRoute,
  indexRoute,
  tenantRoute.addChildren([
    dashboardRoute,
    assetRoutes,
    mutationRoutes,
    // ... other routes
  ]),
]);

export const router = createRouter({ routeTree });
```

### 4.2 Protected Routes

```typescript
// src/routes/protected.tsx
export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user } = useAppSelector((s) => s.auth);
  const navigate = useNavigate();
  
  useEffect(() => {
    if (!user) {
      navigate({ to: "/login" });
    }
  }, [user, navigate]);
  
  if (!user) return null;
  
  return <>{children}</>;
}
```

### 4.3 Permission-Based Rendering

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

---

## 5. COMPONENTS

### 5.1 Component Hierarchy

```
App
├── RouterProvider
│   ├── RootRoute (__root.tsx)
│   │   ├── LoginPage (public)
│   │   └── TenantRoute (protected)
│   │       └── DashboardShell
│   │           ├── Sidebar
│   │           ├── Topbar
│   │           └── KeepAliveOutlet
│   │               ├── RecentTabsBar
│   │               └── [Page Content]
│   │                   ├── PageContainer
│   │                   ├── DataTable
│   │                   └── Forms
└── Toaster (global)
```

### 5.2 Reusable Components

**DataTable:**
```typescript
// src/components/commons/data-table/DataTable.tsx
export function DataTable<TData>({
  columns,
  data,
  isLoading,
  pagination,
  onPaginationChange,
}: DataTableProps<TData>) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    // ... other options
  });
  
  return (
    <div>
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {flexRender(header.column.columnDef.header, header.getContext())}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {isLoading ? (
            <TableSkeleton />
          ) : (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
      <DataTablePagination table={table} />
    </div>
  );
}
```

**PageContainer:**
```typescript
// src/components/commons/containers/PageContainer.tsx
export function PageContainer({
  title,
  description,
  actions,
  children,
}: PageContainerProps) {
  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">{title}</h1>
          {description && <p className="text-muted-foreground">{description}</p>}
        </div>
        {actions && <div className="flex gap-2">{actions}</div>}
      </div>
      {children}
    </div>
  );
}
```

---

## 6. FORMS & VALIDATION

### 6.1 Form with React Hook Form + Zod

```typescript
// Define schema
const assetSchema = z.object({
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  category_id: z.number().min(1, "Category is required"),
  status: z.enum(["ready", "broken", "maintenance"]),
  purchase_date: z.string().optional(),
});

type AssetFormData = z.infer<typeof assetSchema>;

// Component
function AssetForm() {
  const form = useForm<AssetFormData>({
    resolver: zodResolver(assetSchema),
    defaultValues: {
      code: "",
      name: "",
      status: "ready",
    },
  });
  
  const [createAsset] = useCreateAssetMutation();
  
  const onSubmit = async (data: AssetFormData) => {
    try {
      await createAsset(data).unwrap();
      toast.success("Asset created");
      form.reset();
    } catch (error) {
      toast.error("Failed to create asset");
    }
  };
  
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="code"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Code</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        
        <Button type="submit">Create Asset</Button>
      </form>
    </Form>
  );
}
```

---

## 7. CUSTOM HOOKS

### 7.1 usePermission

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
  
  return { hasPermission, hasAnyPermission };
}
```

### 7.2 useDebounce

```typescript
// src/hooks/useDebounce.ts
export function useDebounce<T>(value: T, delay: number = 500): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);
  
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);
    
    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);
  
  return debouncedValue;
}
```

---

## 8. BFF (Backend for Frontend)

### 8.1 Hono Proxy

```typescript
// apps/bff/src/index.ts
import { Hono } from 'hono'
import { cors } from 'hono/cors'

const app = new Hono()
const GO_SERVICE_URL = process.env.GO_SERVICE_URL || 'http://localhost:8080/api'

// Enable CORS
app.use('*', cors({
  origin: ['http://localhost:5173'],
  credentials: true,
}))

// Proxy all requests
app.all('*', async (c) => {
  const url = new URL(c.req.url)
  const targetUrl = `${GO_SERVICE_URL}${url.pathname}${url.search}`
  
  const response = await fetch(targetUrl, {
    method: c.req.method,
    headers: c.req.header(),
    body: ['GET', 'HEAD'].includes(c.req.method) 
      ? undefined 
      : await c.req.arrayBuffer(),
  })
  
  return new Response(response.body, {
    status: response.status,
    headers: response.headers
  })
})

export default app
```

**Purpose:**
- Handle CORS
- Proxy requests to backend
- Future: Add caching, rate limiting, request transformation

---

## 9. STYLING

### 9.1 Tailwind CSS

```typescript
// Utility function
import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// Usage
<div className={cn(
  "rounded-lg border p-4",
  isActive && "bg-primary text-primary-foreground",
  className
)}>
  Content
</div>
```

### 9.2 Shadcn/UI Components

Shadcn/UI adalah collection of reusable components yang bisa di-copy ke project:

```bash
# Add component
npx shadcn-ui@latest add button
npx shadcn-ui@latest add dialog
npx shadcn-ui@latest add table
```

Components tersimpan di `src/components/ui/` dan bisa di-customize.

---

## 10. BEST PRACTICES

### 10.1 DO's

✅ Gunakan TypeScript strict mode
✅ Semua API call via RTK Query
✅ Form validation dengan Zod
✅ Permission check di UI
✅ Responsive design (mobile-first)
✅ Lazy load pages & components
✅ Memoize expensive computations
✅ Use semantic HTML
✅ Accessible components (ARIA)
✅ Error boundaries

### 10.2 DON'Ts

❌ Jangan fetch data di useEffect (gunakan RTK Query)
❌ Jangan simpan sensitive data di localStorage
❌ Jangan bypass permission check
❌ Jangan inline styles (gunakan Tailwind)
❌ Jangan prop drilling (gunakan Redux/Context)
❌ Jangan mutate state directly
❌ Jangan hardcode API URLs
❌ Jangan skip loading & error states

---

## 11. PERFORMANCE OPTIMIZATION

### 11.1 Code Splitting

```typescript
// Lazy load pages
const AssetListPage = lazy(() => import("../modules/asset/pages/AssetListPage"));

// Lazy load components
const AssetTable = lazy(() => import("../layouts/AssetTable"));
```

### 11.2 Virtualization

```typescript
// For large lists
import { useVirtualizer } from "@tanstack/react-virtual";

function VirtualList({ items }) {
  const parentRef = useRef<HTMLDivElement>(null);
  
  const virtualizer = useVirtualizer({
    count: items.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 50,
  });
  
  return (
    <div ref={parentRef} style={{ height: "400px", overflow: "auto" }}>
      <div style={{ height: `${virtualizer.getTotalSize()}px` }}>
        {virtualizer.getVirtualItems().map((virtualItem) => (
          <div
            key={virtualItem.key}
            style={{
              position: "absolute",
              top: 0,
              left: 0,
              width: "100%",
              height: `${virtualItem.size}px`,
              transform: `translateY(${virtualItem.start}px)`,
            }}
          >
            {items[virtualItem.index]}
          </div>
        ))}
      </div>
    </div>
  );
}
```

### 11.3 Memoization

```typescript
// Memoize expensive computations
const filteredAssets = useMemo(() => {
  return assets.filter((asset) => 
    asset.name.toLowerCase().includes(search.toLowerCase())
  );
}, [assets, search]);

// Memoize callbacks
const handleClick = useCallback(() => {
  console.log("Clicked");
}, []);
```

---

**Next:** [Database Schema](./DATABASE_SCHEMA.md)
