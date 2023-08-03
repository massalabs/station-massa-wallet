import * as path from 'path';

import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

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
