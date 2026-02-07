type LogLevel = "debug" | "info" | "warn" | "error";

const prefix = "[{{ package_name | kebab_case }}]";

export function log(level: LogLevel, message: string, meta?: unknown) {
	if (meta === undefined) {
		// biome-ignore lint/suspicious/noConsole: expected for app logger
		console[level](`${prefix} ${message}`);
		return;
	}

	// biome-ignore lint/suspicious/noConsole: expected for app logger
	console[level](`${prefix} ${message}`, meta);
}
