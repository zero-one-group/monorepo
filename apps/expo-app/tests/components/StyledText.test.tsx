import { renderAsync, screen } from "@testing-library/react-native";
import React from "react";

import { MonoText } from "#/components/StyledText";

test("MonoText renders children", async () => {
	await renderAsync(<MonoText>hello</MonoText>);
	expect(screen.getByText("hello")).toBeOnTheScreen();
});
