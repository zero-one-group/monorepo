import React from "react";
import { Alert, Linking, Pressable, PressableProps } from "react-native";

export interface ExternalLinkProps extends PressableProps {
	href: string;
}

export function ExternalLink({
	href,
	...props
}: React.PropsWithChildren<ExternalLinkProps>) {
	const onPress = React.useCallback(async () => {
		const supported = await Linking.canOpenURL(href);

		if (supported) {
			await Linking.openURL(href);
		} else {
			Alert.alert(`Don't know how to open this URL: ${href}`);
		}
	}, [href]);

	return <Pressable onPress={onPress} {...props} />;
}
