import React from "react";
import { renderAsync, screen } from "@testing-library/react-native";

import { MonoText } from "#/components/StyledText";

test("MonoText renders children", async () => {
	await renderAsync(<MonoText>Hello World</MonoText>);
	expect(screen.getByText("Hello World")).toBeOnTheScreen();
});
