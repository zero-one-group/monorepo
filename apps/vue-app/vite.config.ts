import { resolve } from "node:path";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
	envPrefix: "VITE_",
	plugins: [vue(), tailwindcss(), tsconfigPaths()],
	server: { port: 3000, host: false },
	publicDir: resolve("public"),
	optimizeDeps: {
		exclude: ["@repo/shared-ui-vue"],
	},
	build: {
		ssr: false,
		minify: true,
		cssMinify: true,
		chunkSizeWarningLimit: 1024 * 2,
		reportCompressedSize: false,
		emptyOutDir: true,
		manifest: true,
	},
});
