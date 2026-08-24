import vue from "@vitejs/plugin-vue";
import { isCI } from "std-env";
import { defineConfig } from "vitest/config";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
	plugins: [vue(), tsconfigPaths()],
	test: {
		environment: "happy-dom",
		exclude: ["node_modules", "tests-e2e"],
		reporters: isCI ? ["html", "github-actions"] : ["html", "default"],
		include: ["./tests/**/*.{test,spec}.ts"],
		setupFiles: ["./tests/setup-client.ts"],
		outputFile: {
			json: "./tests-results/vitest-results.json",
			html: "./tests-results/index.html",
		},
		coverage: {
			provider: "v8",
			reporter: ["html-spa", "text-summary"],
			reportsDirectory: "./tests-results/coverage",
			cleanOnRerun: true,
			clean: true,
		},
		globals: true,
	},
});
