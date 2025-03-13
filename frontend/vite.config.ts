import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tsconfigPaths from "vite-tsconfig-paths"

export default defineConfig(({ command }) => ({
  plugins: [
    react({
      babel: {
        plugins: [
          ["babel-plugin-react-compiler"],
        ],
      },
    }), 
    tsconfigPaths()
  ],
  ...(command === 'serve' && {
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:8081',
          changeOrigin: true,
          secure: false,
          ws: true, 
        },
      },
      watch: {
        usePolling: true, 
      },
      hmr: {
        overlay: true, 
      },
    },
  }),
}));
