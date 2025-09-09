import "./styles/global.css";
import { Links, Meta, Outlet, Scripts, ScrollRestoration } from "react-router";
import { isRouteErrorResponse } from "react-router";
import NotFound from "#/components/errors/404";
import InternalError from "#/components/errors/500";
import type { Route } from "./+types/root";

export const links: Route.LinksFunction = () => [
	{
		rel: "preconnect",
		href: "https://cdn.jsdelivr.net",
		crossOrigin: "anonymous",
	},
];

export function Layout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="en">
			<head>
				<meta charSet="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
				<Meta />
				<Links />
			</head>
			<body suppressHydrationWarning={true}>
				{children}
				<ScrollRestoration />
				<Scripts />
			</body>
		</html>
	);
}

export default function App() {
	return <Outlet />;
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
	const isError = isRouteErrorResponse(error);

	let message = "Oops!";
	let details = "An unexpected error occurred.";
	let stack: string | undefined;

	if (isError && error.status === 404) {
		return <NotFound />;
	}

	if (isError) {
		message = "Error";
		details = error.statusText || details;
	} else if (import.meta.env.DEV && error && error instanceof Error) {
		details = error.message;
		stack = error.stack;
	}

	return <InternalError message={message} details={details} stack={stack} />;
}
