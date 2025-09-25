import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
	index("routes/home/page.tsx"),
	route("login", "routes/auth/login/page.tsx"),
] satisfies RouteConfig;
