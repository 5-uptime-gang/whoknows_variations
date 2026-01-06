const { defineConfig, devices } = require("@playwright/test");

const baseURL = process.env.BASE_URL || "http://localhost:8080";

module.exports = defineConfig({
  testDir: "./tests",
  timeout: 90 * 1000,
  expect: {
    timeout: 15 * 1000,
  },
  use: {
    baseURL,
    headless: true,
    trace: "retain-on-failure",
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],
});
