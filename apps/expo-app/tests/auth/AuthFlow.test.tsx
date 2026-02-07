import {
	fireEvent,
	renderAsync,
	screen,
	waitFor,
} from "@testing-library/react-native";
import React from "react";

import { DashboardPanel } from "#/features/auth/DashboardPanel";
import { LoginForm } from "#/features/auth/LoginForm";
import { AuthProvider, useAuth } from "#/state/auth";

function AuthExample() {
	const { user } = useAuth();
	return user ? <DashboardPanel /> : <LoginForm />;
}

test("user can sign in then sign out (state flows through context)", async () => {
	await renderAsync(
		<AuthProvider>
			<AuthExample />
		</AuthProvider>,
	);

	fireEvent.changeText(screen.getByLabelText("Email"), "test@example.com");
	fireEvent.changeText(screen.getByLabelText("Password"), "pass1234");
	fireEvent.press(screen.getByRole("button", { name: "Sign in" }));

	await waitFor(() => {
		expect(screen.getByTestId("dashboard-welcome")).toHaveTextContent(
			"Welcome, test@example.com",
		);
	});

	fireEvent.press(screen.getByRole("button", { name: "Sign out" }));
	await waitFor(() => {
		expect(screen.getByRole("button", { name: "Sign in" })).toBeOnTheScreen();
	});
});

test("invalid credentials show a helpful error", async () => {
	await renderAsync(
		<AuthProvider>
			<LoginForm />
		</AuthProvider>,
	);

	fireEvent.changeText(screen.getByLabelText("Email"), "not-an-email");
	fireEvent.changeText(screen.getByLabelText("Password"), "1234");
	fireEvent.press(screen.getByRole("button", { name: "Sign in" }));

	await waitFor(() => {
		expect(screen.getByTestId("auth-error")).toHaveTextContent(
			"Enter a valid email.",
		);
	});
});
