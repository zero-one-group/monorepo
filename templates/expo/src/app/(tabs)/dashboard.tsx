import React from "react";
import { StyleSheet } from "react-native";

import { Redirect } from "expo-router";

import { Screen } from "#/components/layout/Screen";
import { View } from "#/components/Themed";
import { DashboardPanel } from "#/features/auth/DashboardPanel";
import { useAuth } from "#/state/auth";

export default function DashboardScreen() {
	const { user } = useAuth();

	if (!user) {
		return <Redirect href="/login" />;
	}

	return (
		<Screen>
			<View style={styles.container}>
				<DashboardPanel />
			</View>
		</Screen>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		justifyContent: "center",
	},
});
