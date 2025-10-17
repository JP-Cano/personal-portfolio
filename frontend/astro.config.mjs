// @ts-check
import node from "@astrojs/node";
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
    server: {
      allowedHosts: true,
    },
  },

  env: {
    schema: {
      PORTFOLIO_BACKEND_URL: envField.string({
        context: "server",
        access: "secret",
      }),
      PORT: envField.number({
        context: "server",
        access: "public",
        default: DEFAULT_PORT,
      }),
      ENVIRONMENT: envField.string({ context: "server", access: "public" }),
    },
  },

  adapter: node({
    mode: "standalone",
  }),
});
