import js from "@eslint/js";
import eslintPluginAstro from "eslint-plugin-astro";
import importPlugin from "eslint-plugin-import";
import eslintPluginSvelte from "eslint-plugin-svelte";
import globals from "globals";
import tseslint from "typescript-eslint";

export default [
  // Base JavaScript recommended rules
  js.configs.recommended,

  // TypeScript ESLint recommended rules
  ...tseslint.configs.recommended,

  // Astro recommended rules
  ...eslintPluginAstro.configs.recommended,

  // Svelte recommended rules (Svelte 5)
  ...eslintPluginSvelte.configs.recommended,
  {
    plugins: {
      import: importPlugin,
    },
    rules: {
      // Import sorting rules
      "import/order": [
        "error",
        {
          groups: [
            "builtin",
            "external",
            "internal",
            "parent",
            "sibling",
            "index",
          ],
          "newlines-between": "never",
          alphabetize: {
            order: "asc",
            caseInsensitive: true,
          },
        },
      ],
      // TypeScript rules
      "@typescript-eslint/no-unused-vars": [
        "warn",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
        },
      ],
      "@typescript-eslint/no-explicit-any": "warn",

      // General code quality rules
      "no-console": ["warn", { allow: ["warn", "error"] }],
      "prefer-const": "error",
      "no-var": "error",

      // Astro-specific rules
      "astro/no-set-html-directive": "error",
      "astro/no-unused-css-selector": "warn",
    },
  },

  // Use the TypeScript parser inside <script lang="ts"> blocks of Svelte files
  // and expose browser globals (these islands run only on the client).
  {
    files: ["**/*.svelte", "**/*.svelte.ts", "**/*.svelte.js"],
    languageOptions: {
      globals: {
        ...globals.browser,
      },
      parserOptions: {
        parser: tseslint.parser,
      },
    },
    rules: {
      // Svelte 5 runes ($props, $state, $derived) are declared with `let` even
      // when never explicitly reassigned, so prefer-const is a false positive.
      "prefer-const": "off",
    },
  },

  // Ignore patterns
  {
    ignores: ["dist/**", "node_modules/**", ".astro/**"],
  },
];
