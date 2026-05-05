import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/sentinel': {
        target: 'http://localhost',
        changeOrigin: true,
      }
    }
  }
})