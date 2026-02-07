import { createFileRoute } from "@tanstack/react-router";
import { Welcome } from "#/components/home/welcome";

export const Route = createFileRoute("/")({
	component: IndexPage,
});

function IndexPage() {
	return <Welcome />;
}
