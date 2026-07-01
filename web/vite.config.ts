import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import type { IncomingMessage, ServerResponse } from 'node:http'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 这里的 /api 是一个匹配规则
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,             // 允许跨域
        rewrite: (path) => path.replace(/^\/api/, ''),
        configure: (proxy) => {
          proxy.on('proxyRes', (proxyRes: IncomingMessage, req: IncomingMessage, res: ServerResponse) => {
            const isSSE = proxyRes.headers['content-type']?.includes('text/event-stream')
            if (!isSSE) return
            // 避免 dev proxy 缓冲 SSE，确保 status/content 逐段到达前端
            res.setHeader('Cache-Control', 'no-cache, no-transform')
            res.setHeader('X-Accel-Buffering', 'no')
            if (typeof res.flushHeaders === 'function') {
              res.flushHeaders()
            }
          })
        },
      }
    }
  },
  build: {
    outDir: 'dist', // 确保打包输出目录名为 dist
  }
})