import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Welcome } from "#/components/home/welcome";

describe("Welcome", () => {
	it("renders hero and cards", () => {
		render(<Welcome />);

		expect(screen.getByText("TanStack Start Template")).toBeInTheDocument();
		expect(screen.getByText("TanStack Start Hosting")).toBeInTheDocument();
		expect(screen.getByText("TanStack Router v1")).toBeInTheDocument();
	});
});
