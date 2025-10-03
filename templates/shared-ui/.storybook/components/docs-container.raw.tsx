import { DocsContainer as BaseContainer } from "@storybook/addon-docs/blocks";
import * as React from "react";
import { dark, light, listenToColorScheme } from "../themes";

const themes = { light, dark };

/**
 * Switch color scheme based on the global types or system preferences
 */
export const DocsContainer: typeof BaseContainer = ({ children, context }) => {
	const [theme, setTheme] = React.useState<"light" | "dark">("light");

	React.useEffect(() => {
		return listenToColorScheme(context.channel, (theme) => {
			setTheme(theme === "dark" ? "dark" : "light");
		});
	}, [context.channel]);

	React.useEffect(() => {
		document.documentElement.dataset.theme = theme;
	}, [theme]);

	return (
		<BaseContainer context={context} theme={themes[theme]}>
			{children}
		</BaseContainer>
	);
};
