import * as Lucide from "lucide-react";
import { Link } from "#/components/link";

const cards = [
    {
        title: "Zero One Starter Kit",
        description:
            "Launch your next project in minutes with our battle-tested monorepo template and development tools.",
        href: "https://github.com/zero-one-group/monorepo",
    },
    {
        title: "TanStack Start Hosting",
        description:
            "Deploy a full-stack Start app with a simple Vite build + Node server entry.",
        href: "https://tanstack.com/start/latest/docs/framework/react/guide/hosting",
    },
    {
        title: "TanStack Router v1",
        description:
            "Build type-safe routing with file-based routes and automatic route tree generation.",
        href: "https://tanstack.com/router/latest",
    },
];

export function Welcome() {
    return (
        <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
            <div className="rounded-2xl border bg-card p-8 shadow-sm">
                <div className="flex flex-col gap-2">
                    <p className="text-muted-foreground text-sm">TanStack Start Template</p>
                    <h1 className="font-semibold text-3xl tracking-tight">
                        Zero One Starter Kit
                    </h1>
                    <p className="max-w-2xl text-muted-foreground">
                        A minimal, full-stack TanStack Start app template wired with TanStack
                        Router file-based routing, Tailwind v4, tests, and Docker.
                    </p>
                </div>
            </div>

            <div className="mt-10 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
                {cards.map((card) => (
                    <div key={card.title} className="rounded-lg border bg-card p-6 shadow-sm">
                        <h3 className="font-medium text-foreground text-lg">{card.title}</h3>
                        <p className="mt-2 text-muted-foreground">{card.description}</p>
                        <Link
                            href={card.href}
                            className="mt-4 inline-flex items-center text-primary hover:opacity-90"
                            newTab
                        >
                            <span>Learn more</span>
                            <Lucide.ChevronsRight className="ml-1 size-5" />
                        </Link>
                    </div>
                ))}
            </div>
        </main>
    );
}
