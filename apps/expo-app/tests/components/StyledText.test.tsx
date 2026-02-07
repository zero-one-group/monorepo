import React from "react";
import { act, create } from "react-test-renderer";

import { MonoText } from "#/components/StyledText";

test("MonoText renders children", () => {
	let rendered: ReturnType<typeof create> | null = null;

	act(() => {
		rendered = create(<MonoText>hello</MonoText>);
	});

	expect(rendered).not.toBeNull();
	expect((rendered as any).toJSON()).toBeTruthy();
});
