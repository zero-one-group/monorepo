import { isProduction } from "std-env";

/** @type {import('next').NextConfig} */
const nextConfig = {
	output: "standalone",
	cleanDistDir: true,
	reactStrictMode: true,
	poweredByHeader: false,
	productionBrowserSourceMaps: false,
	images: { remotePatterns: [{ protocol: "https", hostname: "**" }] },
	eslint: { ignoreDuringBuilds: isProduction },
	typescript: { ignoreBuildErrors: isProduction },
	logging: { fetches: { fullUrl: true } },
};

export default nextConfig;
