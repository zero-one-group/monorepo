import type { PropsWithChildren } from "react";
import { StyleSheet } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { View } from "#/components/Themed";

type ScreenProps = PropsWithChildren;

export function Screen({ children }: ScreenProps) {
	return (
		<SafeAreaView style={styles.safeArea}>
			<View style={styles.container}>{children}</View>
		</SafeAreaView>
	);
}

const styles = StyleSheet.create({
	safeArea: {
		flex: 1,
	},
	container: {
		flex: 1,
	},
});
