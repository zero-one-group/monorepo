/**
 * Reference: https://playwright.dev/docs/test-configuration
 */

import "dotenv/config";
import { resolve } from "node:path";
import { defineConfig, devices } from "@playwright/test";
import { env, isCI } from "std-env";

export const STORAGE_STATE = resolve(".playwright/user.json");
const APP_PORT = env.PORT || 3300;

export default defineConfig({
	quiet: !!isCI,
	testDir: resolve("tests-e2e"),
	outputDir: resolve("tests-results/e2e"),
	fullyParallel: true,
	forbidOnly: !!isCI,
	retries: isCI ? 2 : 0,
	workers: isCI ? 1 : undefined,
	reporter: [
		["html", { open: "never", outputFolder: "./tests-results/e2e-html" }],
		["list"],
	],
	preserveOutput: "always",
	use: {
		...devices["Desktop Chrome"],
		baseURL: env.URL || `http://localhost:${APP_PORT}`,
		defaultBrowserType: "chromium",
		colorScheme: "no-preference",
		ignoreHTTPSErrors: true,
		locale: "en-US",
		trace: "on-all-retries",
		video: { mode: isCI ? "off" : "on" },
		contextOptions: {
			recordVideo: { dir: resolve("tests-results/e2e-videos") },
			reducedMotion: "reduce",
		},
		screenshot: isCI ? "off" : "only-on-failure",
		launchOptions: {
			slowMo: 1000,
		},
	},
	projects: [
		{ name: "Chromium", use: { ...devices["Desktop Chrome"] } },
		{ name: "Firefox", use: { ...devices["Desktop Firefox"] } },
		{ name: "Safari", use: { ...devices["Desktop Safari"] } },
		{ name: "Mobile Chrome", use: { ...devices["Pixel 5"] } },
	],
	webServer: env.URL
		? undefined
		: {
				command: "moon tanstack-start:build && moon tanstack-start:start",
				reuseExistingServer: !isCI,
				port: Number(APP_PORT),
				timeout: 30_000,
			},
});
