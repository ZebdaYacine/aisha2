import { defineConfig, devices } from "@playwright/test";

const externalBaseURL = process.env.PLAYWRIGHT_BASE_URL;
const baseURL = externalBaseURL ?? "http://localhost:3100";

export default defineConfig({
  testDir: "./tests/e2e",
  fullyParallel: false,
  workers: 1,
  use: { baseURL, trace: "on-first-retry" },
  webServer: externalBaseURL ? undefined : [
    { command: "node tests/e2e/catalogue-fixture.mjs", url: "http://127.0.0.1:4080/api/v1/categories", reuseExistingServer: false },
    { command: "API_BASE_URL=http://127.0.0.1:4080/api/v1 npm run dev -- --webpack --port 3100", url: "http://localhost:3100/api/health", reuseExistingServer: false },
  ],
  projects: [
    { name: "desktop", use: { ...devices["Desktop Chrome"] } },
    { name: "mobile", use: { ...devices["Pixel 7"] } },
  ],
});
