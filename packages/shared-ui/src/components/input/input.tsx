import * as Lucide from "lucide-react";
import * as React from "react";
import { Button } from "../button/button";
import { toast } from "../toast/toast";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "../tooltip/tooltip";
import { inputStyles } from "./input.css";

export interface InputProps
	extends React.InputHTMLAttributes<HTMLInputElement> {
	onCopy?: () => void;
	showCopyButton?: boolean;
	showExternalCopyButton?: boolean;
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
	(
		{
			className,
			type,
			showCopyButton,
			showExternalCopyButton,
			onCopy,
			...props
		},
		ref,
	) => {
		const [showPassword, setShowPassword] = React.useState(false);
		const styles = inputStyles({
			hasPasswordToggle: type === "password",
			hasCopyButton: showCopyButton && !type,
		});

		const togglePassword = () => setShowPassword(!showPassword);

		const handleCopy = () => {
			if (!props.value) {
				toast.error("Nothing to copy");
				return;
			}
			navigator.clipboard.writeText(props.value.toString());
			toast.success("Copied to clipboard");
			onCopy?.();
		};

		const PasswordToggle = () => (
			<TooltipProvider>
				<Tooltip delayDuration={150}>
					<TooltipTrigger asChild>
						<button
							type="button"
							onClick={togglePassword}
							className={styles.toggleButton()}
						>
							{showPassword ? (
								<Lucide.EyeOff className="size-4" strokeWidth={2} />
							) : (
								<Lucide.Eye className="size-4" strokeWidth={2} />
							)}
						</button>
					</TooltipTrigger>
					<TooltipContent side="top">
						<p>{showPassword ? "Hide password" : "Show password"}</p>
					</TooltipContent>
				</Tooltip>
			</TooltipProvider>
		);

		const CopyButton = () => (
			<TooltipProvider>
				<Tooltip>
					<TooltipTrigger asChild>
						<button
							type="button"
							onClick={handleCopy}
							className={styles.copyButton()}
						>
							<Lucide.Copy className="size-4" strokeWidth={2} />
						</button>
					</TooltipTrigger>
					<TooltipContent side="top">
						<p>Copy to clipboard</p>
					</TooltipContent>
				</Tooltip>
			</TooltipProvider>
		);

		const ExternalCopyButton = () => (
			<TooltipProvider>
				<Tooltip>
					<TooltipTrigger asChild>
						<Button
							type="button"
							variant="outline"
							size="icon"
							onClick={handleCopy}
							className="shrink-0"
						>
							<Lucide.Copy className="size-4" strokeWidth={2} />
						</Button>
					</TooltipTrigger>
					<TooltipContent side="top">
						<p>Copy to clipboard</p>
					</TooltipContent>
				</Tooltip>
			</TooltipProvider>
		);

		if (showExternalCopyButton) {
			return (
				<div className={styles.wrapperWithCopy()}>
					<div className={styles.wrapper()}>
						<input
							type={showPassword ? "text" : type}
							className={styles.input({ className })}
							ref={ref}
							{...props}
						/>
						{type === "password" && <PasswordToggle />}
					</div>
					<ExternalCopyButton />
				</div>
			);
		}

		return (
			<div className={styles.wrapper()}>
				<input
					type={showPassword ? "text" : type}
					className={styles.input({ className })}
					ref={ref}
					{...props}
				/>
				{type === "password" && <PasswordToggle />}
				{showCopyButton && !type && <CopyButton />}
			</div>
		);
	},
);

Input.displayName = "Input";

export { Input };
