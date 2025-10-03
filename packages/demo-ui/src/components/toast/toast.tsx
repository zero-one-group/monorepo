import { Toaster as Sonner, toast } from "sonner";
import { useTheme } from "../../theme";
import { toastStyles } from "./toast.css";

type ToasterProps = React.ComponentProps<typeof Sonner>;

const Toaster = ({ ...props }: ToasterProps) => {
	const { theme = "system" } = useTheme();
	const styles = toastStyles();

	return (
		<Sonner
			theme={theme as ToasterProps["theme"]}
			className="toaster group"
			toastOptions={{
				classNames: {
					toast: styles.toast(),
					description: styles.description(),
					actionButton: styles.actionButton(),
					cancelButton: styles.cancelButton(),
				},
			}}
			{...props}
		/>
	);
};

export { Toaster, toast };
