export const SCOPES = {
  // Master Data & Lokasi
  ROOM: {
    CREATE: "room:create",
    READ: "room:read",
    UPDATE: "room:update",
    DELETE: "room:delete",
  },
  BED: {
    CREATE: "bed:create",
    READ: "bed:read",
    UPDATE: "bed:update",
    DELETE: "bed:delete",
  },

  // Atribut Asset
  CATEGORY: {
    CREATE: "category:create",
    READ: "category:read",
    UPDATE: "category:update",
    DELETE: "category:delete",
  },
  VENDOR: {
    CREATE: "vendor:create",
    READ: "vendor:read",
    UPDATE: "vendor:update",
    DELETE: "vendor:delete",
  },
  BRAND: {
    CREATE: "brand:create",
    READ: "brand:read",
    UPDATE: "brand:update",
    DELETE: "brand:delete",
  },
  MODEL: {
    CREATE: "model:create",
    READ: "model:read",
    UPDATE: "model:update",
    DELETE: "model:delete",
  },
  ASSET_STATUS: {
    CREATE: "asset_status:create",
    READ: "asset_status:read",
    UPDATE: "asset_status:update",
    DELETE: "asset_status:delete",
  },

  // Modul Asset Utama
  ASSET: {
    CREATE: "asset:create",
    READ: "asset:read",
    UPDATE: "asset:update",
    DELETE: "asset:delete",
  },
  ASSET_MUTATION: {
    CREATE: "asset_mutation:create",
    READ: "asset_mutation:read",
  },
  STOCK_OPNAME: {
    CREATE: "stock_opname:create",
    READ: "stock_opname:read",
    UPDATE: "stock_opname:update",
    CLOSE: "stock_opname:close",
  },

  // Pemeliharaan & Keluhan
  MAINTENANCE: {
    CREATE: "maintenance:create",
    READ: "maintenance:read",
    UPDATE: "maintenance:update",
    DELETE: "maintenance:delete",
    COMPLETE: "maintenance:complete",
  },
  MAINTENANCE_DOCUMENT: {
    UPLOAD: "maintenance_document:upload",
    READ: "maintenance_document:read",
    DELETE: "maintenance_document:delete",
  },
  COMPLAINT: {
    CREATE: "complaint:create",
    READ: "complaint:read",
    UPDATE: "complaint:update",
    DELETE: "complaint:delete",
  },

  // Dokumen & Laporan
  DOCUMENT: {
    CREATE: "document:create",
    READ: "document:read",
    UPDATE: "document:update",
    DELETE: "document:delete",
  },
  REPORT: {
    EXPORT: "report:export",
  },

  // System & Management
  NOTIFICATION: {
    SEND: "notification:send",
  },
  PERMISSION: {
    READ: "permission:read",
  },
  ROLE: {
    CREATE: "role:create",
    READ: "role:read",
    READ_PERMISSIONS: "role:read_permissions",
    UPDATE: "role:update",
    DELETE: "role:delete",
    ASSIGN: "role:assign_permissions",
    REVOKE: "role:revoke_permissions",
  },
  USER: {
    CREATE: "user:create",
    READ: "user:read",
    UPDATE: "user:update",
    DELETE: "user:delete",
    RESET_PWD: "user:reset_password",
  },
  AUDIT: {
    READ: "audit:read",
  },
  TENANT: {
    READ: "tenant:read",
    CREATE: "tenant:create",
    UPDATE: "tenant:update",
    UPDATE_STATUS: "tenant:update_status",
  },
} as const;

// Helper Type untuk TypeScript (Opsional)
export type ScopeType = typeof SCOPES;