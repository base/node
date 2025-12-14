```js
#!/usr/bin/env node
```
/**
 * Simple sanity-check for the `.env.example` file.
 *
 * This script:
 * - verifies that `.env.example` exists;
 * - prints out non-empty keys;
 * - warns if there are obviously placeholder values.
 *
 * Usage:
 *   node scripts/check-env-example.js
 */
 
const fs = require("fs");
const path = require("path");

const ENV_EXAMPLE_PATH = path.join(process.cwd(), ".env.example");

function main() {
  if (!fs.existsSync(ENV_EXAMPLE_PATH)) {
    console.error("❌ .env.example not found in project root.");
    process.exitCode = 1;
    return;
  }

  const content = fs.readFileSync(ENV_EXAMPLE_PATH, "utf8");
  const lines = content.split(/\r?\n/);

  const entries = [];
  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith("#")) continue;

    const [key, ...rest] = trimmed.split("=");
    const value = rest.join("=").trim();

    entries.push({ key: key.trim(), value });
  }

  if (entries.length === 0) {
    console.warn("⚠️ .env.example is empty or only contains comments.");
    return;
  }
    console.log("🔎 Checking .env.example variables:");
  for (const { key, value } of entries) {
    const isPlaceholder =
      value === "" ||
      value.toLowerCase().includes("your_") ||
      value.toLowerCase().includes("changeme");

    const status = isPlaceholder ? "placeholder" : "ok";
    console.log(`- ${key}: ${status}`);
  }

  console.log("\n✅ Check finished. Make sure to keep real secrets in .env only.");
}

main();
