FROM oven/bun:1-alpine

WORKDIR /app

# Copy package.json
COPY package.json ./

# Install dependencies
RUN bun install

# Copy source code
COPY src ./src
COPY tsconfig.json ./

# Expose port
EXPOSE 3000

# Run
CMD ["bun", "run", "src/index.ts"]
