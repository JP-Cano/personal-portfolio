// @ts-check
import node from "@astrojs/node";
import svelte from "@astrojs/svelte";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig, envField } from "astro/config";

const DEFAULT_PORT = 4321;

// https://astro.build/config
export default defineConfig({
  output: "server",

  server: {
    host: true,
    allowedHosts: true,
  },

  vite: {
    assetsInclude: ["**/*.docx"],

    server: {
      allowedHosts: true,
    },

    plugins: [tailwindcss()],
  },

  env: {
    schema: {
      PORTFOLIO_BACKEND_URL: envField.string({
        context: "server",
        access: "secret",
        default: "http://localhost:8080/api/v1",
      }),
      PORT: envField.number({
        context: "server",
        access: "public",
        default: DEFAULT_PORT,
      }),
      ENVIRONMENT: envField.string({
        context: "server",
        access: "public",
        default: "development",
      }),
    },
  },

  adapter: node({
    mode: "standalone",
    trustProxy: true,
  }),

  integrations: [svelte()],
});
