/** @type {import("prettier").Config} */
export default {
  // Astro plugin for formatting .astro files
  plugins: ["prettier-plugin-astro"],

  // General formatting rules
  semi: true,
  singleQuote: false,
  tabWidth: 2,
  useTabs: false,
  trailingComma: "es5",
  printWidth: 80,
  bracketSpacing: true,
  arrowParens: "always",
  endOfLine: "lf",

  // Astro-specific overrides
  overrides: [
    {
      files: "*.astro",
      options: {
        parser: "astro",
      },
    },
  ],
};
