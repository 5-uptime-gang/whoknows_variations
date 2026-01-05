const { test, expect } = require("@playwright/test");
const {
  uniqueUser,
  registerUser,
  loginUser,
  logoutUser,
} = require("./helpers");

test("user can register and sees authenticated navbar", async ({ page }) => {
  const user = uniqueUser();
  await registerUser(page, user);

  await expect(page.locator("#nav-logout")).toBeVisible();
  await expect(page.locator("#nav-login")).toBeHidden();
  await expect(page.locator("#nav-register")).toBeHidden();
});

test("logout clears session and login works again", async ({ page }) => {
  const user = uniqueUser();
  await registerUser(page, user);
  await logoutUser(page);
  await loginUser(page, user);
});
