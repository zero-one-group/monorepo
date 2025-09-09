import { GLOBALS_UPDATED } from "@storybook/core-events";
import { create } from "@storybook/theming";

const brand = {
	brandTitle: "UI Components",
	brandUrl: "/",
	fontBase:
		'"Inter", Helvetica, Arial, system-ui, sans-serif, "Segoe UI", Roboto',
	appBorderRadius: 4,
	inputBorderRadius: 4,
};

export const light = create({
	base: "light",
	...brand,
	colorPrimary: "#0f131a",
	colorSecondary: "#0f131a",

	// UI
	appBg: "#ffffff",
	appContentBg: "#ffffff",
	appPreviewBg: "#ffffff",
	appBorderColor: "#e5e7eb",

	// Text colors
	textColor: "#0f131a",
	textMutedColor: "#6b7280",
	textInverseColor: "#ffffff",

	// Toolbar colors
	barTextColor: "#6b7280",
	barHoverColor: "#374151",
	barSelectedColor: "#1fa2ff",
	barBg: "#ffffff",

	// Form colors
	buttonBg: "#f3f4f6",
	buttonBorder: "#e5e7eb",
	inputBg: "#ffffff",
	inputBorder: "#e5e7eb",
	inputTextColor: "#0f131a",
});

export const dark = create({
	base: "dark",
	...brand,
	// brandImage: '/images/logo-name-dark.png',
	colorPrimary: "#ffffff",
	colorSecondary: "#374151",

	// UI
	appBg: "#0f131a",
	appContentBg: "#0f131a",
	appPreviewBg: "#0f131a",
	appBorderColor: "#374151",

	// Text colors
	textColor: "#ffffff",
	textMutedColor: "#9ca3af",
	textInverseColor: "#0f131a",

	// Toolbar colors
	barTextColor: "#9ca3af",
	barHoverColor: "#d1d5db",
	barSelectedColor: "#1fa2ff",
	barBg: "#0f131a",

	// Form colors
	buttonBg: "#1f2937",
	buttonBorder: "#374151",
	inputBg: "#111827",
	inputBorder: "#374151",
	inputTextColor: "#ffffff",
});

export type Theme = "light" | "dark" | "system";

type EventListener = (
	eventName: string,
	callback: (context: { globals: Record<string, unknown> }) => void,
) => void;

export function listenToColorScheme(
	eventEmitter: { on: EventListener; off: EventListener },
	callback: (theme: Theme) => void,
) {
	const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");

	const handleMediaChange = (event: MediaQueryListEvent) => {
		callback(event.matches ? "dark" : "light");
	};

	const handleGlobalsChange = ({ globals }) => {
		const theme = globals.theme as Theme;

		if (theme === "system") {
			callback(mediaQuery.matches ? "dark" : "light");
			mediaQuery.addEventListener("change", handleMediaChange);
		} else {
			callback(theme);
			mediaQuery.removeEventListener("change", handleMediaChange);
		}
	};

	const initColorScheme = () => {
		const globals = new URL(window.location.href).searchParams.get("globals");

		if (globals) {
			const [key, value] = globals.split(":");
			if (key === "theme") {
				return handleGlobalsChange({ globals: { theme: value } });
			}
		}

		handleGlobalsChange({ globals: { theme: "system" } });
	};

	initColorScheme();

	eventEmitter.on(GLOBALS_UPDATED, handleGlobalsChange);

	return () => {
		eventEmitter.off(GLOBALS_UPDATED, handleGlobalsChange);
	};
}
