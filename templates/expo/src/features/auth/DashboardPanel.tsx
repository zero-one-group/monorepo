import React from "react";
import { Pressable, StyleSheet } from "react-native";

import { Text, View } from "#/components/Themed";
import { useAuth } from "#/state/auth";

export function DashboardPanel() {
	const { user, signOut } = useAuth();

	if (!user) return null;

	return (
		<View style={styles.container}>
			<Text accessibilityRole="header" style={styles.title}>
				Dashboard
			</Text>
			<Text testID="dashboard-welcome" style={styles.body}>
				Welcome, {user.email}
			</Text>

			<Pressable
				accessibilityRole="button"
				accessibilityLabel="Sign out"
				onPress={signOut}
				style={({ pressed }) => [
					styles.button,
					{ opacity: pressed ? 0.75 : 1 },
				]}
			>
				<Text style={styles.buttonText}>Sign out</Text>
			</Pressable>
		</View>
	);
}

const styles = StyleSheet.create({
	container: {
		paddingHorizontal: 20,
		paddingVertical: 16,
	},
	title: {
		fontSize: 22,
		fontWeight: "700",
	},
	body: {
		marginTop: 10,
		fontSize: 15,
		opacity: 0.85,
	},
	button: {
		marginTop: 18,
		alignSelf: "flex-start",
		borderRadius: 10,
		paddingHorizontal: 12,
		paddingVertical: 10,
		backgroundColor: "#e74c3c",
	},
	buttonText: {
		color: "#fff",
		fontWeight: "700",
		fontSize: 14,
	},
});
