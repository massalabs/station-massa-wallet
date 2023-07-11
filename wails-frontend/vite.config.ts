import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import * as path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      events: 'rollup-plugin-node-polyfills/polyfills/events',
      '@': path.resolve(__dirname, 'src'),
      '@wailsjs': path.resolve(__dirname, 'wailsjs'),
    },
  },
  build: {
    emptyOutDir: true,
  },
});
