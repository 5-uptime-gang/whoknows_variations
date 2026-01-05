const { expect } = require("@playwright/test");

function uniqueUser() {
  const slug = `${Date.now()}-${Math.random().toString(16).slice(2, 6)}`;
  return {
    username: `pw-user-${slug}`,
    email: `pw-user-${slug}@example.com`,
    password: "Pw!123456",
  };
}

async function registerUser(page, user) {
  await page.goto("/register");
  await page.fill('input[name="username"]', user.username);
  await page.fill('input[name="email"]', user.email);
  await page.fill('input[name="password"]', user.password);
  await page.fill('input[name="password2"]', user.password);
  await Promise.all([page.waitForURL("**/"), page.click('button[type="submit"]')]);
  await expect(page.locator("#nav-logout")).toBeVisible({ timeout: 10_000 });
}

async function logoutUser(page) {
  const logout = page.locator("#nav-logout");
  await expect(logout).toBeVisible({ timeout: 10_000 });
  await logout.click();
  await page.waitForLoadState("networkidle");
  await expect(page.locator("#nav-login")).toBeVisible({ timeout: 10_000 });
}

async function loginUser(page, user) {
  await page.goto("/login");
  await page.fill('input[name="username"]', user.username);
  await page.fill('input[name="password"]', user.password);
  await Promise.all([page.waitForURL("**/"), page.click("#login-button")]);
  await expect(page.locator("#nav-logout")).toBeVisible({ timeout: 10_000 });
}

async function ensureLoggedIn(page, user) {
  await registerUser(page, user);
  // Navbar JS might not refresh state immediately; hard reload helps ensure cookie is picked up
  await page.reload({ waitUntil: "networkidle" });
  await expect(page.locator("#nav-logout")).toBeVisible({ timeout: 10_000 });
}

module.exports = {
  uniqueUser,
  registerUser,
  loginUser,
  logoutUser,
  ensureLoggedIn,
};
