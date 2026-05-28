# Build stage
FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Copy root package.json
COPY package.json bun.lock* ./

# Copy workspace packages
COPY apps/frontend/package.json ./apps/frontend/
COPY apps/src ./apps/src

# Install all dependencies
RUN bun install

# Copy frontend source
COPY apps/frontend ./apps/frontend

# Build frontend
WORKDIR /app/apps/frontend
RUN bun run build

# Production stage
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/apps/frontend/dist /usr/share/nginx/html

# Copy nginx config
COPY apps/frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
