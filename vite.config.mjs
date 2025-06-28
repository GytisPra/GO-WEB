import { defineConfig } from 'vite'
import path from 'path'

export default defineConfig({
  root: 'web',
  build: {
    outDir: 'static',
    emptyOutDir: false,
    rollupOptions: {
      input: {
        main: path.resolve('web/ts/main.ts'),
      },
      output: {
        entryFileNames: 'js/[name].js',
        assetFileNames: 'css/[name][extname]',
      },
    },
  },
  css: {
    postcss: path.resolve('./postcss.config.mjs'),
  },
})
