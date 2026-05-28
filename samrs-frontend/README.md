# SAMRS Frontend

React + TypeScript SPA for the SAMRS hospital asset management platform.

## Tech Stack

- **React 19** + TypeScript
- **Vite** (Rolldown)
- **Redux Toolkit** + RTK Query (state & data fetching)
- **TanStack Router** + TanStack Table
- **shadcn/ui** + Tailwind CSS
- **React Hook Form** + Zod (forms & validation)
- **Bun** (package manager & runtime)

## Structure

```
apps/
├── frontend/   # React SPA — feature-modular
│   └── src/
│       ├── modules/    # Feature modules (asset, complaint, maintenance, ...)
│       ├── components/ # Shared UI components
│       ├── routes/     # TanStack Router routes
│       ├── store/      # Redux store + RTK Query API
│       └── hooks/      # Shared hooks
└── bff/        # Backend-for-Frontend (Node.js/Bun)
```

## Getting Started

```bash
bun install
bun run dev
```

## Build

```bash
bun run build
```

See the root [README](../README.md) for full project setup with Docker.
