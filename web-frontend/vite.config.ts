import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '', // in index.html, assets files start with ./ so that they are relative to the html file
  build: {
    outDir: '../internal/handler/html/dist',
    emptyOutDir: true,
    manifest: true,
    sourcemap: true,
    assetsDir: './', // put the assets next to the index.html file
  },
});
