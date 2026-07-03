type LogLevel = "debug" | "info" | "warn" | "error";

const prefix = "[{{ package_name | kebab_case }}]";

export function log(level: LogLevel, message: string, meta?: unknown) {
	if (meta === undefined) {
		// oxlint-disable-next-line no-console
		console[level](`${prefix} ${message}`);
		return;
	}

	// oxlint-disable-next-line no-console
	console[level](`${prefix} ${message}`, meta);
}
