import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: './swagger.json',
  output: {
    path: './src/api',
    format: 'prettier',
  },
  plugins: [
    '@hey-api/client-axios',
    '@tanstack/react-query',
    '@hey-api/typescript',
  ],
});
