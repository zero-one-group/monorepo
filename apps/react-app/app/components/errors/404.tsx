import { useNavigate } from "react-router";
import { Link } from "#/components/link";
import { errorStyles } from "./error.css";

export default function NotFound() {
	const navigate = useNavigate();
	const styles = errorStyles();

	const handleBack = () => {
		if (window.history.length > 1) {
			navigate(-1);
		} else {
			navigate("/");
		}
	};

	return (
		<div className={styles.wrapper()}>
			<div className={styles.decorativeGradient()}>
				<div className={styles.gradientInner()}>
					<div className={styles.gradientBg()} />
				</div>
			</div>
			<div className={styles.decorativeCode()}>
				<h2 className={styles.decorativeText()}>404</h2>
			</div>
			<div className={styles.content()}>
				<div className={styles.container()}>
					<p className={styles.errorCode()}>404</p>
					<h1 className={styles.title()}>Page not found</h1>
					<p className={styles.description()}>
						Sorry, we couldn't find the page you're looking for.
					</p>
					<div className={styles.actions()}>
						<button
							type="button"
							onClick={handleBack}
							className={styles.primaryButton()}
						>
							Go back
						</button>
						<Link href="#" className={styles.secondaryButton()}>
							Docs
						</Link>
					</div>
				</div>
			</div>
		</div>
	);
}
