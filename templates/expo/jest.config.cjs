const fs = require("node:fs");
const path = require("node:path");

/**
 * Expo + pnpm workspaces can result in a nested React copy under
 * @testing-library/react-native, which breaks Hooks in Jest.
 *
 * This mapper forces a single React + renderer instance.
 */
function resolveFromNodeModules(pkg) {
	const local = path.join(__dirname, "node_modules", pkg);
	if (fs.existsSync(local)) return local;
	return path.join(__dirname, "..", "..", "node_modules", pkg);
}

/** @type {import('jest').Config} */
module.exports = {
	preset: "jest-expo",
	setupFilesAfterEnv: ["<rootDir>/jest.setup.ts"],
	moduleNameMapper: {
		"^#/(.*)$": "<rootDir>/src/$1",
		"^react$": resolveFromNodeModules("react"),
	},
	testPathIgnorePatterns: ["/node_modules/", "/dist/", "/web-build/"],
	transformIgnorePatterns: [
		"node_modules/(?!(?:.pnpm/)?((jest-)?react-native|@react-native(-community)?|expo(nent)?|@expo(nent)?/.*|@expo-google-fonts/.*|react-navigation|@react-navigation/.*|@sentry/react-native|native-base|react-native-svg))",
	],
};
