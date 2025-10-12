import js from "@eslint/js";
import eslintPluginAstro from "eslint-plugin-astro";
import importPlugin from "eslint-plugin-import";
import tseslint from "typescript-eslint";

export default [
  // Base JavaScript recommended rules
  js.configs.recommended,

  // TypeScript ESLint recommended rules
  ...tseslint.configs.recommended,

  // Astro recommended rules
  ...eslintPluginAstro.configs.recommended,
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

  // Ignore patterns
  {
    ignores: ["dist/**", "node_modules/**", ".astro/**"],
  },
];
