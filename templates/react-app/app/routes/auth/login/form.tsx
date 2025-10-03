import { useActionState } from "react";
import { useFormStatus } from "react-dom";
import { Link } from "#/components/link";
import { type LoginState, loginAction } from "./action";

function SubmitButton() {
	const { pending } = useFormStatus();

	return (
		<button
			type="submit"
			className="w-full rounded-lg bg-blue-600 px-4 py-2 text-white font-semibold hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
			disabled={pending}
		>
			{pending ? "Signing in..." : "Sign in"}
		</button>
	);
}

export function LoginForm() {
	const [state, formAction] = useActionState<LoginState, FormData>(
		loginAction,
		{},
	);

	return (
		<form className="mt-8 space-y-6" action={formAction}>
			{state?.error && (
				<div className="text-center text-red-500 text-sm">{state.error}</div>
			)}
			{state?.data && (
				<div className="text-center text-green-600 text-sm">{`Hello ${state.data.email}`}</div>
			)}

			<div className="space-y-4">
				<div>
					<label htmlFor="email" className="sr-only">
						Email address
					</label>
					<input
						id="email"
						name="email"
						type="email"
						className="relative block w-full appearance-none rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-blue-500 focus:outline-none focus:ring-blue-500 sm:text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
						placeholder="Email address"
					/>
				</div>
				<div>
					<label htmlFor="password" className="sr-only">
						Password
					</label>
					<input
						id="password"
						name="password"
						type="password"
						className="relative block w-full appearance-none rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:border-blue-500 focus:outline-none focus:ring-blue-500 sm:text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
						placeholder="Password"
					/>
				</div>
			</div>

			<div className="flex items-center justify-between">
				<div className="flex items-center">
					<input
						id="remember-me"
						name="remember-me"
						type="checkbox"
						className="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
					/>
					<label
						htmlFor="remember-me"
						className="ml-2 block text-gray-900 text-sm dark:text-gray-300"
					>
						Remember me
					</label>
				</div>

				<Link
					href="#"
					className="font-medium text-blue-600 text-sm hover:text-blue-500 dark:text-blue-400"
				>
					Forgot your password?
				</Link>
			</div>

			<SubmitButton />
		</form>
	);
}
