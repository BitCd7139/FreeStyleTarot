import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 这里的 /api 是一个匹配规则
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,             // 允许跨域
      }
    }
  }
})