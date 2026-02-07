const path = require("path");
const { getDefaultConfig } = require("expo/metro-config");

const projectRoot = __dirname;
const monorepoRoot = path.resolve(projectRoot, "../..");

/** @type {import('expo/metro-config').MetroConfig} */
const config = getDefaultConfig(projectRoot);

// Watch all files in the monorepo so workspace packages resolve correctly.
config.watchFolders = [monorepoRoot];

// Resolve modules from the app first, then from the monorepo root.
config.resolver.nodeModulesPaths = [
	path.resolve(projectRoot, "node_modules"),
	path.resolve(monorepoRoot, "node_modules"),
];

// Prevent Metro from walking up directories to resolve modules.
config.resolver.disableHierarchicalLookup = true;

module.exports = config;
