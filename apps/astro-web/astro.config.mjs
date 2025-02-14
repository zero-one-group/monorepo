import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'astro/config'

export default defineConfig({
  server: { port: 4321, host: false },
  vite: {
    plugins: [tailwindcss()],
  },
})
