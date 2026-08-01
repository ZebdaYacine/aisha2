import { expect, test } from "@playwright/test";
test("browses catalogue and opens a product", async ({ page }) => {
  await page.goto("/en/products");
  await expect(
    page.getByRole("heading", { name: "Products", exact: true }),
  ).toBeVisible();
  await page
    .getByRole("link", { name: "Kabyle silver brooch", exact: true })
    .last()
    .click();
  await expect(
    page.getByRole("heading", { name: "Kabyle silver brooch", level: 1 }),
  ).toBeVisible();
});
test("Arabic storefront is RTL", async ({ page }) => {
  await page.goto("/ar");
  await expect(page.locator("html")).toHaveAttribute("dir", "rtl");
  await expect(page.getByRole("heading", { level: 1 })).toBeVisible();
});
test("adds an item to cart", async ({ page }) => {
  const pageErrors: string[] = [];
  page.on("pageerror", (error) => pageErrors.push(error.message));
  await page.goto("/en/products/kabyle-silver-brooch");
  await page.waitForTimeout(1000);
  await page.getByTestId("product-add-to-cart").click();
  expect(pageErrors).toEqual([]);
  await expect(page.getByRole("dialog", { name: "Cart" })).toBeVisible();
});
