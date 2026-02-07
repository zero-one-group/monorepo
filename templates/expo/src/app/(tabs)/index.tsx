import { StyleSheet } from "react-native";
import { Screen } from "#/components/layout/Screen";
import { Text, View } from "#/components/Themed";

export default function HomeScreen() {
	return (
		<Screen>
			<View style={styles.container}>
				<Text style={styles.title}>Expo + React Native</Text>
				<Text style={styles.subtitle}>
					Template skeleton with Expo Router, Jest, and Biome.
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
