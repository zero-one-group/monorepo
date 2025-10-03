import type { Decorator } from "@storybook/react-vite";
import * as React from "react";
import { Toaster } from "../../src/components";
import { ThemeProvider } from "../../src/theme";

function setTheme(theme: "light" | "dark") {
	document.documentElement.dataset.theme = theme;
}

export const withThemeProvider: Decorator = (Story, context) => {
	const theme = context.parameters.theme || context.globals.theme;

	React.useEffect(() => {
		if (theme !== "system") {
			setTheme(theme);
			return undefined;
		}

		const query = window.matchMedia("(prefers-color-scheme: dark)");

		setTheme(query.matches ? "dark" : "light");

		const handleChange = (event: MediaQueryListEvent) => {
			setTheme(event.matches ? "dark" : "light");
		};

		query.addEventListener("change", handleChange);

		return () => {
			query.removeEventListener("change", handleChange);
		};
	}, [theme]);

	return (
		<ThemeProvider>
			{/* @ts-ignore */}
			<Story />
			<Toaster />
		</ThemeProvider>
	);
};
