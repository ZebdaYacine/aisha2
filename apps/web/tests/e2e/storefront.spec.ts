import { expect, test } from "@playwright/test";

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => localStorage.setItem("aisha:splash-seen-at", String(Date.now())));
});

test("browses catalogue and opens a product", async ({ page }) => {
  await page.goto("/en/products");
  await expect(
    page.getByRole("heading", { name: "Products", exact: true }),
  ).toBeVisible();
  const productLink = page
    .getByRole("link", { name: "Kabyle silver brooch", exact: true })
    .first();
  await expect(productLink).toHaveAttribute(
    "href",
    "/en/products/kabyle-silver-brooch",
  );
  await Promise.all([
    page.waitForURL(/\/en\/products\/kabyle-silver-brooch/),
    productLink.click(),
  ]);
  await expect(
    page.getByRole("heading", { name: "Kabyle silver brooch", level: 1 }),
  ).toBeVisible();
});
test("Arabic storefront is RTL", async ({ page }) => {
  await page.goto("/ar");
  await expect(page.locator("html")).toHaveAttribute("dir", "rtl");
  await expect(page.getByRole("heading", { level: 1 })).toBeVisible();
});

test("registration form follows the locale-specific panel placement", async ({ page }) => {
  const formPanel = page.locator("main > div > section");

  await page.goto("/en/register");
  await expect(formPanel).toHaveAttribute("dir", "ltr");
  await expect(formPanel).toHaveClass(/lg:order-1/);

  await page.goto("/ar/register");
  await expect(formPanel).toHaveAttribute("dir", "rtl");
  await expect(formPanel).toHaveClass(/lg:order-2/);
});

test("Arabic public catalogue pages preserve RTL and localized controls", async ({ page }, testInfo) => {
  await page.goto("/ar/products");
  await expect(page.locator("html")).toHaveAttribute("dir", "rtl");
  if (testInfo.project.name === "mobile") {
    await expect(page.getByRole("button", { name: "الفلاتر" })).toBeVisible();
  } else {
    await expect(page.getByText("الترتيب", { exact: true })).toBeVisible();
  }
  await page.goto("/ar/artisans");
  await expect(page.getByRole("heading", { level: 1 })).toContainText("تعرّف على الصنّاع");
});
test("public storefront remains usable at a mobile viewport", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/en");
  await expect(page.getByRole("heading", { level: 1 })).toBeVisible();
  await expect(page.getByRole("button", { name: "Open menu" })).toBeVisible();
  await page.goto("/en/products");
  await expect(page.getByRole("button", { name: "Filters" })).toBeVisible();
});
test("does not invent product availability before inventory is implemented", async ({ page }) => {
  const pageErrors: string[] = [];
  page.on("pageerror", (error) => pageErrors.push(error.message));
  await page.goto("/en/products/kabyle-silver-brooch");
  await page.waitForTimeout(1000);
  await expect(page.getByTestId("product-add-to-cart")).toBeDisabled();
  expect(pageErrors).toEqual([]);
});

test("login maps rate limiting and protects account navigation", async ({ page }) => {
  await page.route("**/api/auth/login", async (route) => route.fulfill({ status: 429, contentType: "application/json", body: JSON.stringify({ error: { code: "RATE_LIMITED" } }) }));
  await page.goto("/en/login");
  await page.getByLabel("Email or phone number").fill("customer@example.com");
  await page.getByLabel("Password").fill("a-secure-password");
  await page.getByRole("button", { name: "Sign in" }).click();
  await expect(page.locator("#auth-login-errors")).toContainText("Too many attempts. Please wait and try again.");
  await page.context().clearCookies();
  await page.goto("/en/account");
  await expect(page).toHaveURL(/\/en\/login/);
});

test("successful login enters the protected account shell", async ({ page }) => {
  await page.route("**/api/auth/login", async (route) => route.fulfill({ status: 200, headers: { "set-cookie": "aisha_access=test-access; Path=/; HttpOnly; SameSite=Lax" }, contentType: "application/json", body: JSON.stringify({ user: { id: "user-1" } }) }));
  await page.goto("/en/login");
  await page.getByLabel("Email or phone number").fill("customer@example.com");
  await page.getByLabel("Password").fill("a-secure-password");
  await page.getByRole("button", { name: "Sign in" }).click();
  await expect(page).toHaveURL(/\/en\/account$/);
  await expect(page.locator("#main-content")).toBeVisible();
});

test("authenticated customer can create an address through the real form contract", async ({ page, context }) => {
  await context.addCookies([{ name: "aisha_access", value: "test-access", domain: "localhost", path: "/", httpOnly: true, sameSite: "Lax" }]);
  let created = false;
  await page.route("**/api/customer/addresses", async (route) => {
    if (route.request().method() === "POST") { created = true; await route.fulfill({ status: 201, contentType: "application/json", body: JSON.stringify({ id: "address-1", ...route.request().postDataJSON() }) }); return; }
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(created ? [{ id: "address-1", fullName: "Amina", phone: "", line1: "12 Rue", line2: "", city: "Alger", postalCode: "16000", country: "Algeria", isDefault: true }] : []) });
  });
  await page.goto("/en/account/addresses");
  for (const [label, value] of [["Full name", "Amina"], ["Address", "12 Rue"], ["City", "Alger"], ["Postal code", "16000"], ["Country", "Algeria"]]) await page.getByLabel(label, { exact: true }).fill(value);
  await page.getByRole("button", { name: "Add Delivery address" }).click();
  await expect(page.getByText("12 Rue")).toBeVisible();
});
