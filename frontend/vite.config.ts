import path from "path"
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/v1": {
        target: process.env.VITE_API_BASE_URL || "http://localhost:9000",
        changeOrigin: true,
        secure: false,
      },
      "/api": {
        target: process.env.VITE_API_BASE_URL || "http://localhost:9000",
        changeOrigin: true,
        secure: false,
      }
    },
  },
})
