const { test, expect } = require("@playwright/test");

test("weather page shows forecast", async ({ page }) => {
  await page.goto("/weather");

  const weatherDay = page.locator("#weather-container .weather-day").first();
  await expect(weatherDay).toBeVisible({ timeout: 20_000 });
  await expect(page.locator("#weather-container")).not.toContainText(
    /Error loading weather data/i
  );
});
