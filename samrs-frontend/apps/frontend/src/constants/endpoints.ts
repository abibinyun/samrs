// src/constants/endpoints.ts

export const API_ENDPOINTS = {
  // 1. AUTHENTICATION & CONNECTIVITY
  AUTH: {
    LOGIN: "/api/v1/auth/login",
    ME: "/api/v1/auth/me",
    PING: "/ping",
    SECURE_PING: "/api/v1/secure-ping",
  },

  // 2. ROOM MANAGEMENT
  ROOMS: {
    BASE: "/api/v1/rooms",
    DETAIL: (id: string) => `/api/v1/rooms/${id}`,
  },

  // 3. BED MANAGEMENT
  BEDS: {
    BASE: "/api/v1/beds",
    DETAIL: (id: string) => `/api/v1/beds/${id}`,
    BY_ROOM: (roomId: string) => `/api/v1/beds?room_id=${roomId}`,
  },

  // 4. ASSET ATTRIBUTES (Kategori, Vendor, Brand, Model, Status)
  ATTRIBUTES: {
    CATEGORIES: {
      BASE: "/api/v1/categories",
      DETAIL: (id: string) => `/api/v1/categories/${id}`,
    },
    VENDORS: {
      BASE: "/api/v1/vendors",
      DETAIL: (id: string) => `/api/v1/vendors/${id}`,
    },
    BRANDS: {
      BASE: "/api/v1/brands",
      DETAIL: (id: string) => `/api/v1/brands/${id}`,
    },
    MODELS: {
      BASE: "/api/v1/models",
      DETAIL: (id: string) => `/api/v1/models/${id}`,
    },
    STATUSES: {
      BASE: "/api/v1/asset-statuses",
      DETAIL: (id: string) => `/api/v1/asset-statuses/${id}`,
    },
  },

  // 5. ASSET INVENTORY
  ASSETS: {
    BASE: "/api/v1/assets",
    DETAIL: (id: string) => `/api/v1/assets/${id}`,
    TIMELINE: (id: string) => `/api/v1/assets/${id}/timeline`,
    QR: (id: string, size: number = 256) => `/api/v1/assets/${id}/qr?size=${size}`,
    PUBLIC_SCAN: (tenantSlug: string, assetCode: string) => 
      `/public/tenants/${tenantSlug}/assets/${assetCode}`,
  },

  // 5.5 COMPLAINT / WORK ORDER
  COMPLAINTS: {
    BASE: "/api/v1/complaints",
    DETAIL: (id: string) => `/api/v1/complaints/${id}`,
  },

  // 5.6 MAINTENANCE & CALIBRATION
  MAINTENANCE: {
    SCHEDULES: {
      BASE: "/api/v1/maintenance-schedules",
      DETAIL: (id: string) => `/api/v1/maintenance-schedules/${id}`,
      COMPLETE: (id: string) => `/api/v1/maintenance-schedules/${id}/complete`,
    },
    DOCUMENTS: {
      UPLOAD: (scheduleId: string) => `/api/v1/maintenance-schedules/${scheduleId}/documents`,
      LIST: (scheduleId: string) => `/api/v1/maintenance-schedules/${scheduleId}/documents`,
      DOWNLOAD: (docId: string) => `/api/v1/maintenance-documents/${docId}/download`,
      DELETE: (docId: string) => `/api/v1/maintenance-documents/${docId}`,
    },
  },

  // 5.9 ASSET MUTATION
  MUTATIONS: {
    BASE: "/api/v1/asset-mutations",
    DETAIL: (id: string) => `/api/v1/asset-mutations/${id}`,
  },

  // 5.10 STOCK OPNAME
  STOCK_OPNAME: {
    BASE: "/api/v1/stock-opnames",
    ITEMS: (opnameId: string) => `/api/v1/stock-opnames/${opnameId}/items`,
    CLOSE: (opnameId: string) => `/api/v1/stock-opnames/${opnameId}/close`,
  },

  // 5.7 REPORT EXPORT
  REPORTS: {
    ASSETS: "/api/v1/reports/assets/export",
    COMPLAINTS: "/api/v1/reports/complaints/export",
    MAINTENANCE: "/api/v1/reports/maintenance-schedules/export",
  },

  // 5.8 NOTIFICATIONS
  NOTIFICATIONS: {
    SEND: "/api/v1/notifications/send",
  },

  // 6. RBAC MANAGEMENT
  RBAC: {
    PERMISSIONS: "/api/v1/permissions",
    ROLES: {
      BASE: "/api/v1/roles",
      DETAIL: (id: string) => `/api/v1/roles/${id}`,
      SET_TENANT_ADMIN: (id: string) => `/api/v1/roles/${id}/tenant-admin`,
      PERMISSIONS: (id: string) => `/api/v1/roles/${id}/permissions`,
    },
    USERS: {
      BASE: "/api/v1/users",
      DETAIL: (id: string) => `/api/v1/users/${id}`,
      UPDATE_ROLE: (id: string) => `/api/v1/users/${id}/role`,
      UPDATE_STATUS: (id: string) => `/api/v1/users/${id}/status`,
      RESET_PASSWORD: (id: string) => `/api/v1/users/${id}/password`,
    },
  },

  // 8. AUDIT TRAIL
  AUDIT: {
    BASE: "/api/v1/audit-trails",
    DETAIL: (id: string) => `/api/v1/audit-trails/${id}`,
  },

  // 9. ADMIN TENANTS (Super Admin Only)
  ADMIN: {
    TENANTS: {
      BASE: "/api/v1/admin/tenants",
      DETAIL: (id: string) => `/api/v1/admin/tenants/${id}`,
      UPDATE_STATUS: (id: string) => `/api/v1/admin/tenants/${id}/status`,
    },
  },

  // 10. GLOBAL DOCUMENT MANAGEMENT
  DOCUMENTS: {
    BASE: "/api/v1/documents",
    DETAIL: (id: string) => `/api/v1/documents/${id}`,
    FILES: (id: string) => `/api/v1/documents/${id}/files`,
    DOWNLOAD_FILE: (fileId: string) => `/api/v1/document-files/${fileId}/download`,
    DELETE_FILE: (fileId: string) => `/api/v1/document-files/${fileId}`,
  },
} as const;