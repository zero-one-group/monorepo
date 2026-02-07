import type { ErrorComponentProps } from "@tanstack/react-router";
import { ErrorComponent } from "@tanstack/react-router";

export function DefaultCatchBoundary({ error, reset }: ErrorComponentProps) {
	return (
		<div className="mx-auto max-w-3xl px-4 py-16">
			<h1 className="font-semibold text-3xl tracking-tight">
				Something went wrong
			</h1>
			<p className="mt-2 text-muted-foreground">
				An unexpected error occurred.
			</p>
			<div className="mt-6 rounded-lg border bg-card p-4">
				<ErrorComponent error={error} />
			</div>
			<button
				type="button"
				onClick={reset}
				className="mt-6 inline-flex items-center rounded-md bg-primary px-4 py-2 font-medium text-primary-foreground"
			>
				Try again
			</button>
		</div>
	);
}
