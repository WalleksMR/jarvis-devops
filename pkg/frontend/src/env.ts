import { z } from "zod";

const envSchema = z.object({
  VITE_API_URL: z.string(),
  VITE_ENABLE_API_DELAY: z.string().transform((value) => value === "true"),
  VITE_BASE_URL: z.string().default("/"),
});

export const env = envSchema.parse(import.meta.env);
