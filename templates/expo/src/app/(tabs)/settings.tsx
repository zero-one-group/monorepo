import { StyleSheet } from "react-native";
import { Screen } from "#/components/layout/Screen";
import { Text, View } from "#/components/Themed";

export default function SettingsScreen() {
	return (
		<Screen>
			<View style={styles.container}>
				<Text style={styles.title}>Settings</Text>
				<Text style={styles.subtitle}>
					This page is a starter for app preferences.
				</Text>
			</View>
		</Screen>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		alignItems: "center",
		justifyContent: "center",
		paddingHorizontal: 20,
	},
	title: {
		fontSize: 20,
		fontWeight: "bold",
	},
	subtitle: {
		marginTop: 12,
		fontSize: 14,
		opacity: 0.75,
		textAlign: "center",
	},
});
