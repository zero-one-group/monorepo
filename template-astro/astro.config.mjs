import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "astro/config";

export default defineConfig({
	server: { port: {{ port_number }}, host: false },
	vite: {
		plugins: [tailwindcss()],
	},
});
