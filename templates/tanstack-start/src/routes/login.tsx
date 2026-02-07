import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/login")({
    component: LoginPage,
});

function LoginPage() {
    const [email, setEmail] = React.useState("");

    return (
        <main className="mx-auto max-w-md px-4 py-12 sm:px-6 lg:px-8">
            <h1 className="font-semibold text-3xl tracking-tight">Sign in</h1>
            <p className="mt-2 text-muted-foreground">
                Demo form only (no auth wired).
            </p>

            <form className="mt-8 rounded-xl border bg-card p-6 shadow-sm">
                <label className="block text-sm">
                    <span className="text-muted-foreground">Email</span>
                    <input
                        value={email}
                        onChange={(e) => setEmail(e.currentTarget.value)}
                        type="email"
                        placeholder="you@company.com"
                        className="mt-2 w-full rounded-md border bg-background px-3 py-2"
                    />
                </label>

                <button
                    type="button"
                    className="mt-4 inline-flex w-full items-center justify-center rounded-md bg-primary px-4 py-2 font-medium text-primary-foreground"
                    onClick={() => alert(`Hello, ${email || "friend"}!`)}
                >
                    Continue
                </button>
            </form>
        </main>
    );
}
