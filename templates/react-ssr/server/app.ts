import "react-router";
import { createRequestHandler } from "@react-router/express";
import express from "express";
import helmet from "helmet";
import { generateCspDirectives } from "./helmet";

declare module "react-router" {
	interface AppLoadContext {
		VALUE_FROM_EXPRESS: string;
	}
}

export const app = express();

app.use(
	helmet({
		xPoweredBy: false,
		contentSecurityPolicy: {
			directives: generateCspDirectives(),
		},
	}),
);

app.use(
	createRequestHandler({
		build: () => import("virtual:react-router/server-build"),
		getLoadContext() {
			return {
				VALUE_FROM_EXPRESS: "Hello from Express",
			};
		},
	}),
);
