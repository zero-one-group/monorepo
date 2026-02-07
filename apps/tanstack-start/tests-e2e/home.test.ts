import { expect, test } from "@playwright/test";

test.describe.configure({ mode: "serial" });

test.describe("@feature-default", () => {
	test("has title", async ({ page }) => {
		await page.goto("/");
		await expect(page).toHaveTitle(/TanStack Start App/);
	});

	test("contain hero title", async ({ page }) => {
		await page.goto("/");
		await expect(
			page.getByRole("heading", { name: "Zero One Starter Kit" }),
		).toBeVisible();
	});
});
