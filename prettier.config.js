/** @type {import('prettier').Config & import('prettier-plugin-tailwindcss').PluginOptions} */
const config = {
  plugins: ["prettier-plugin-tailwindcss"],
  htmlWhitespaceSensitivity: "ignore",
  printWidth: 100,
  bracketSameLine: false,
};

export default config;
