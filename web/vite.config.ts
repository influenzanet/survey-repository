import { defineConfig } from 'vite';
export default defineConfig({
  server: {
    proxy: {
      "^/namespace":"http://127.0.0.1:8080",
      "^/survey":"http://127.0.0.1:8080",
    },
    headers: {
    },
  },
});