import path from "node:path";
import { existsSync } from "node:fs";

import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  console.log(`\x1b[34mLoading environment variables for mode: "${mode}"`);
  let envFile = `.env`;

  if (!existsSync(envFile)) {
    console.warn(
      `\x1b[33m⚠️  O arquivo de variáveis de ambiente "${envFile}" não foi encontrado.\n\x1b[0m`,
    );
  }

  return {
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
        "@public": path.resolve(__dirname, "./public"),
      },
    },
    mode,
    server: {
      port: Number(process.env.VITE_PORT) || 3000,
    },
    base: process.env.VITE_BASE_URL || "/",
  };
});
