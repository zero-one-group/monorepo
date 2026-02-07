import { Link } from "expo-router";
import { Pressable, StyleSheet } from "react-native";
import { Screen } from "#/components/layout/Screen";
import { Text, View } from "#/components/Themed";
import { useAuth } from "#/state/auth";

export default function HomeScreen() {
	const { user, signOut } = useAuth();

	return (
		<Screen>
			<View style={styles.container}>
				<Text style={styles.title}>Expo + React Native</Text>
				<Text style={styles.subtitle}>
					Template skeleton with Expo Router, Jest, and Biome.
				</Text>
				<Text style={styles.status}>
					{user ? `Signed in as ${user.email}` : "Not signed in"}
				</Text>

				<View style={styles.actions}>
					{user ? (
						<>
							<Link href="/dashboard" asChild>
								<Pressable style={styles.primaryButton}>
									<Text style={styles.primaryButtonText}>Go to dashboard</Text>
								</Pressable>
							</Link>
							<Pressable
								accessibilityRole="button"
								accessibilityLabel="Sign out"
								onPress={signOut}
								style={styles.secondaryButton}
							>
								<Text style={styles.secondaryButtonText}>Sign out</Text>
							</Pressable>
						</>
					) : (
						<>
							<Link href="/login" asChild>
								<Pressable style={styles.primaryButton}>
									<Text style={styles.primaryButtonText}>Open login</Text>
								</Pressable>
							</Link>
							<Link href="/dashboard" asChild>
								<Pressable style={styles.secondaryButton}>
									<Text style={styles.secondaryButtonText}>
										Try dashboard
									</Text>
								</Pressable>
							</Link>
						</>
					)}
				</View>
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
	status: {
		marginTop: 18,
		fontSize: 13,
		opacity: 0.7,
	},
	actions: {
		marginTop: 16,
		gap: 10,
		width: "100%",
		maxWidth: 320,
	},
	primaryButton: {
		paddingVertical: 12,
		borderRadius: 12,
		backgroundColor: "#2ecc71",
		alignItems: "center",
	},
	primaryButtonText: {
		color: "#fff",
		fontWeight: "700",
	},
	secondaryButton: {
		paddingVertical: 12,
		borderRadius: 12,
		backgroundColor: "#ecf0f1",
		alignItems: "center",
	},
	secondaryButtonText: {
		color: "#2c3e50",
		fontWeight: "700",
	},
});
