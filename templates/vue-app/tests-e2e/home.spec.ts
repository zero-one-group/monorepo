import { expect, test } from "@playwright/test";

test("home page has a heading", async ({ page }) => {
	await page.goto("/");
	await expect(page.locator("h1")).toBeVisible();
});

test("404 page shows 404", async ({ page }) => {
	await page.goto("/non-existent-route");
	await expect(page.locator("h1")).toContainText("404");
});
