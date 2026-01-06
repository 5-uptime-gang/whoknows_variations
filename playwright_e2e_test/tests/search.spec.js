const { test, expect } = require("@playwright/test");
const { uniqueUser, ensureLoggedIn } = require("./helpers");

test("search returns results for Docker", async ({ page }) => {
  await ensureLoggedIn(page, uniqueUser());

  await page.goto("/");
  await page.fill("#search-input", "Docker");
  await page.click("#search-button");

  const firstResult = page.locator("#results .search-result-title").first();
  await expect(firstResult).toBeVisible({ timeout: 15_000 });
  await expect(firstResult).toContainText(/docker/i);
});
