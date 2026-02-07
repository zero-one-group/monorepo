import React from "react";
import { StyleSheet } from "react-native";

import { Redirect, useRouter } from "expo-router";

import { Screen } from "#/components/layout/Screen";
import { Text, View } from "#/components/Themed";
import { LoginForm } from "#/features/auth/LoginForm";
import { useAuth } from "#/state/auth";

export default function LoginScreen() {
	const router = useRouter();
	const { user } = useAuth();

	if (user) {
		return <Redirect href="/dashboard" />;
	}

	return (
		<Screen>
			<View style={styles.container}>
				<LoginForm onSignedIn={() => router.replace("/dashboard")} />
				<Text style={styles.hint}>
					Tip: use any email and a 4+ character password.
				</Text>
			</View>
		</Screen>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
	},
	hint: {
		paddingHorizontal: 20,
		paddingBottom: 20,
		fontSize: 13,
		opacity: 0.65,
		textAlign: "center",
	},
});
