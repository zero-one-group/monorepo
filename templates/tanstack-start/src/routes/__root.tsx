/// <reference types="vite/client" />

import * as React from "react";
import {
    HeadContent,
    Link,
    Outlet,
    Scripts,
    createRootRoute,
} from "@tanstack/react-router";
import type { ErrorComponentProps } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { DefaultCatchBoundary } from "#/components/errors/DefaultCatchBoundary";
import { NotFound } from "#/components/errors/NotFound";

import appCss from "../styles/global.css?url";

export const Route = createRootRoute({
    head: () => ({
        meta: [
            { charSet: "utf-8" },
            { name: "viewport", content: "width=device-width, initial-scale=1" },
            { title: "TanStack Start App" },
        ],
        links: [
            { rel: "stylesheet", href: appCss },
            { rel: "icon", href: "/favicon.svg", type: "image/svg+xml" },
        ],
    }),
    errorComponent: (props: ErrorComponentProps) => (
        <RootDocument>
            <DefaultCatchBoundary {...props} />
        </RootDocument>
    ),
    notFoundComponent: () => (
        <RootDocument>
            <NotFound />
        </RootDocument>
    ),
    component: RootComponent,
});

function RootComponent() {
    return (
        <RootDocument>
            <Outlet />
        </RootDocument>
    );
}

function RootDocument({ children }: { children: React.ReactNode }) {
    return (
        <html lang="en">
            <head>
                <HeadContent />
            </head>
            <body>
                <div className="min-h-screen bg-background text-foreground">
                    <header className="border-b bg-card">
                        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
                            <div className="flex items-center gap-3">
                                <img
                                    src="/images/logo-light.svg"
                                    alt="Logo"
                                    className="block h-8 w-auto dark:hidden"
                                />
                                <img
                                    src="/images/logo-dark.svg"
                                    alt="Logo"
                                    className="hidden h-8 w-auto dark:block"
                                />
                            </div>
                            <nav className="flex items-center gap-5 text-sm">
                                <Link
                                    to="/"
                                    activeOptions={ { exact: true } }
                                    activeProps={ { className: "text-foreground" } }
                                    className="text-muted-foreground hover:text-foreground"
                                >
                                    Home
                                </Link>
                                <Link
                                    to="/dashboard"
                                    activeProps={ { className: "text-foreground" } }
                                    className="text-muted-foreground hover:text-foreground"
                                >
                                    Dashboard
                                </Link>
                                <Link
                                    to="/login"
                                    activeProps={ { className: "text-foreground" } }
                                    className="text-muted-foreground hover:text-foreground"
                                >
                                    Sign In
                                </Link>
                                <a
                                    href="/api/health"
                                    className="text-muted-foreground hover:text-foreground"
                                >
                                    Health
                                </a>
                            </nav>
                        </div>
                    </header>

                    {children}
                </div>

                {import.meta.env.DEV ? (
                    <TanStackRouterDevtools position="bottom-right" />
                ) : null}
                <Scripts />
            </body>
        </html>
    );
}
