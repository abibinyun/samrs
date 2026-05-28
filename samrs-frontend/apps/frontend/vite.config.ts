import path from "path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { visualizer } from "rollup-plugin-visualizer";

export default defineConfig({
  // 🔒 KUNCI ROOT KE apps/frontend
  root: path.resolve(__dirname),

  plugins: [
    react(),
    tailwindcss(),
    visualizer({
      open: true,
      gzipSize: true,
      brotliSize: true,
    }),
  ],

  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },

  server: {
    proxy: {
      "/api": {
        target: "http://localhost:3000",
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api/, ""),
      },
    },
  },

  build: {
    // 🔒 KUNCI OUTPUT
    outDir: path.resolve(__dirname, "dist"),
    emptyOutDir: true,

    rollupOptions: {
      output: {
        advancedChunks: {
          groups: [
            {
              name: "react-vendor",
              test: /\/node_modules\/(react|react-dom|scheduler)\//,
            },
            { name: "tanstack", test: /\/node_modules\/@tanstack\// },
            { name: "radix", test: /\/node_modules\/@radix-ui\// },
            { name: "i18n", test: /\/node_modules\/(i18next|react-i18next)\// },
            { name: "date", test: /\/node_modules\/date-fns\// },
            { name: "lodash", test: /\/node_modules\/lodash\// },

            { name: "asset", test: /\/src\/modules\/asset\// },
            { name: "room", test: /\/src\/modules\/room\// },
            { name: "notification", test: /\/src\/modules\/notification\// },
            { name: "mutation", test: /\/src\/modules\/mutation\// },

            { name: "vendor", test: /\/node_modules\// },
          ],
        },
      },
    },
  },
});
