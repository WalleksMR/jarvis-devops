import path from "node:path";
import { existsSync } from "node:fs";

import { config } from "dotenv";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  console.log(`\x1b[34mLoading environment variables for mode: "${mode}"`);

  let envDir = undefined;
  let envFile = `.env`;

  if (!existsSync(envFile)) {
    console.warn(
      `\x1b[33m⚠️  O arquivo de variáveis de ambiente "${envFile}" não foi encontrado.\n\x1b[0m`,
    );
  }

  config({ path: envFile });

  return {
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    mode,
    envDir,
  };
});
