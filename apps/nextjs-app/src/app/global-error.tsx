"use client"; // Error boundaries must be Client Components

import { useRouter } from "next/navigation";
import { errorStyles } from "#/styles/error.css";

interface GlobalErrorProps {
	error: Error & { digest?: string };
	reset: () => void;
}

export default function GlobalError({ error, reset }: GlobalErrorProps) {
	const router = useRouter();
	const styles = errorStyles();

	const handleBack = () => {
		if (window.history.length > 1) {
			router.back();
		} else {
			router.push("/");
		}
	};

	return (
		<html lang="en">
			<body>
				<div className={styles.wrapper()}>
					<div className={styles.decorativeGradient()}>
						<div className={styles.gradientInner()}>
							<div className={styles.gradientBg()} />
						</div>
					</div>
					<div className={styles.decorativeCode()}>
						<h2 className={styles.decorativeText()}>500</h2>
					</div>
					<div className={styles.content()}>
						<div className={styles.container()}>
							<p className={styles.errorCode()}>{error.digest || 500}</p>
							<h1 className={styles.title()}>Internal Server Error</h1>
							<p className={styles.description()}>{error.message}</p>
							{error.stack && (
								<pre className={styles.pre()}>
									<code className={styles.code()}>{error.stack}</code>
								</pre>
							)}
							<div className={styles.actions()}>
								<button
									type="button"
									onClick={() => reset()}
									className={styles.primaryButton()}
								>
									Try again
								</button>
								<button
									type="button"
									onClick={handleBack}
									className={styles.secondaryButton()}
								>
									Go back
								</button>
							</div>
						</div>
					</div>
				</div>
			</body>
		</html>
	);
}
