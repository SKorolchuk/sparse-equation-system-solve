import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  // Strict base URL mapping for seamless deployment to github pages
  // NOTE: Replace 'network-solver' with your actual repository name if it differs
  base: '/network-solver/',
  build: {
    // Vite compiles production assets into a clean layout one level above React project context
    outDir: '../dist',
    emptyOutDir: true,
  },
})
