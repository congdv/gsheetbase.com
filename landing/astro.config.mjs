import { defineConfig } from 'astro/config';
import tailwindcss from '@tailwindcss/vite';
import cloudflare from '@astrojs/cloudflare';

export default defineConfig({
  site: 'https://gsheetbase.com',
  output: 'server',   // or 'hybrid' if most pages are static

  vite: {
    plugins: [tailwindcss()],
  },

  adapter: cloudflare(),
});