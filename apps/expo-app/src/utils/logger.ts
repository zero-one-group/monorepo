type LogLevel = "debug" | "info" | "warn" | "error";

const prefix = "[expo-app]";

export function log(level: LogLevel, message: string, meta?: unknown) {
	if (meta === undefined) {
		// biome-ignore lint/suspicious/noConsole: expected for app logger
		console[level](`${prefix} ${message}`);
		return;
	}

	// biome-ignore lint/suspicious/noConsole: expected for app logger
	console[level](`${prefix} ${message}`, meta);
}
