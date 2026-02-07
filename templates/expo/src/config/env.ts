type EnvKey = `EXPO_PUBLIC_${string}`;

function getEnvString(key: EnvKey, fallback = ""): string {
	const value = process.env[key];
	return typeof value === "string" && value.length > 0 ? value : fallback;
}

export const env = {
	apiBaseUrl: getEnvString("EXPO_PUBLIC_API_BASE_URL", "https://example.com"),
} as const;
