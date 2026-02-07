import React, { useState } from "react";
import {
	ActivityIndicator,
	Pressable,
	StyleSheet,
	TextInput,
} from "react-native";

import { Text, useThemeColor, View } from "#/components/Themed";
import { useAuth } from "#/state/auth";

type LoginFormProps = {
	onSignedIn?: () => void;
};

export function LoginForm({ onSignedIn }: LoginFormProps) {
	const { signIn, status, error, clearError } = useAuth();
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	const isLoading = status === "signingIn";
	const textColor = useThemeColor({}, "text");
	const tint = useThemeColor({}, "tint");
	const border = useThemeColor({}, "tint");
	const background = useThemeColor({}, "background");

	async function handleSignIn() {
		try {
			await signIn({ email, password });
			onSignedIn?.();
		} catch {
			// Error is handled in AuthProvider state.
		}
	}

	return (
		<View style={styles.container}>
			<Text accessibilityRole="header" style={styles.title}>
				Login
			</Text>
			<Text style={styles.subtitle}>
				Example auth flow using a shared context store.
			</Text>

			<View style={styles.form}>
				<Text style={styles.label}>Email</Text>
				<TextInput
					accessibilityLabel="Email"
					autoCapitalize="none"
					autoCorrect={false}
					keyboardType="email-address"
					placeholder="name@company.com"
					placeholderTextColor={String(textColor) + "80"}
					style={[
						styles.input,
						{ borderColor: border, color: textColor, backgroundColor: background },
					]}
					value={email}
					onChangeText={(value) => {
						setEmail(value);
						if (error) clearError();
					}}
					editable={!isLoading}
				/>

				<Text style={[styles.label, styles.labelSpacing]}>Password</Text>
				<TextInput
					accessibilityLabel="Password"
					autoCapitalize="none"
					autoCorrect={false}
					placeholder="••••"
					placeholderTextColor={String(textColor) + "80"}
					secureTextEntry
					style={[
						styles.input,
						{ borderColor: border, color: textColor, backgroundColor: background },
					]}
					value={password}
					onChangeText={(value) => {
						setPassword(value);
						if (error) clearError();
					}}
					editable={!isLoading}
				/>

				{error ? (
					<Text testID="auth-error" style={styles.error}>
						{error}
					</Text>
				) : null}

				<Pressable
					accessibilityRole="button"
					accessibilityLabel="Sign in"
					onPress={handleSignIn}
					disabled={isLoading}
					style={({ pressed }) => [
						styles.button,
						{
							backgroundColor: tint,
							opacity: pressed ? 0.85 : 1,
						},
					]}
				>
					{isLoading ? (
						<ActivityIndicator color="#fff" />
					) : (
						<Text style={styles.buttonText}>Sign in</Text>
					)}
				</Pressable>
			</View>
		</View>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		justifyContent: "center",
		paddingHorizontal: 20,
	},
	title: {
		fontSize: 22,
		fontWeight: "700",
		textAlign: "center",
	},
	subtitle: {
		marginTop: 10,
		fontSize: 14,
		opacity: 0.75,
		textAlign: "center",
	},
	form: {
		marginTop: 24,
	},
	label: {
		fontSize: 13,
		fontWeight: "600",
		opacity: 0.8,
	},
	labelSpacing: {
		marginTop: 12,
	},
	input: {
		marginTop: 8,
		borderWidth: 1,
		borderRadius: 10,
		paddingHorizontal: 12,
		paddingVertical: 10,
		fontSize: 16,
	},
	error: {
		marginTop: 12,
		fontSize: 13,
		color: "#c0392b",
	},
	button: {
		marginTop: 16,
		borderRadius: 12,
		paddingVertical: 12,
		alignItems: "center",
		justifyContent: "center",
	},
	buttonText: {
		color: "#fff",
		fontWeight: "700",
		fontSize: 16,
	},
});
