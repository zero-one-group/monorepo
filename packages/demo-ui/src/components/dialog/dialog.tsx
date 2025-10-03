import * as Lucide from "lucide-react";
import { Dialog as DialogPrimitive } from "radix-ui";
import * as React from "react";
import { dialogStyles } from "./dialog.css";

const Dialog = DialogPrimitive.Root;
const DialogTrigger = DialogPrimitive.Trigger;
const DialogPortal = DialogPrimitive.Portal;
const DialogClose = DialogPrimitive.Close;

const DialogOverlay = React.forwardRef<
	React.ComponentRef<typeof DialogPrimitive.Overlay>,
	React.ComponentPropsWithoutRef<typeof DialogPrimitive.Overlay>
>(({ className, ...props }, ref) => {
	const styles = dialogStyles();
	return (
		<DialogPrimitive.Overlay
			ref={ref}
			className={styles.overlay({ className })}
			{...props}
		/>
	);
});

const DialogContent = React.forwardRef<
	React.ComponentRef<typeof DialogPrimitive.Content>,
	React.ComponentPropsWithoutRef<typeof DialogPrimitive.Content>
>(({ className, children, ...props }, ref) => {
	const styles = dialogStyles();
	return (
		<DialogPortal>
			<DialogOverlay />
			<DialogPrimitive.Content
				ref={ref}
				className={styles.content({ className })}
				{...props}
			>
				{children}
				<DialogPrimitive.Close className={styles.close()}>
					<Lucide.XIcon className={styles.closeIcon()} strokeWidth={2} />
					<span className="sr-only">Close</span>
				</DialogPrimitive.Close>
			</DialogPrimitive.Content>
		</DialogPortal>
	);
});

const DialogHeader = ({
	className,
	...props
}: React.HTMLAttributes<HTMLDivElement>) => {
	const styles = dialogStyles();
	return <div className={styles.header({ className })} {...props} />;
};

const DialogFooter = ({
	className,
	...props
}: React.HTMLAttributes<HTMLDivElement>) => {
	const styles = dialogStyles();
	return <div className={styles.footer({ className })} {...props} />;
};

const DialogTitle = React.forwardRef<
	React.ComponentRef<typeof DialogPrimitive.Title>,
	React.ComponentPropsWithoutRef<typeof DialogPrimitive.Title>
>(({ className, ...props }, ref) => {
	const styles = dialogStyles();
	return (
		<DialogPrimitive.Title
			ref={ref}
			className={styles.title({ className })}
			{...props}
		/>
	);
});

const DialogDescription = React.forwardRef<
	React.ComponentRef<typeof DialogPrimitive.Description>,
	React.ComponentPropsWithoutRef<typeof DialogPrimitive.Description>
>(({ className, ...props }, ref) => {
	const styles = dialogStyles();
	return (
		<DialogPrimitive.Description
			ref={ref}
			className={styles.description({ className })}
			{...props}
		/>
	);
});

Dialog.displayName = "Dialog";
DialogTrigger.displayName = "DialogTrigger";
DialogPortal.displayName = "DialogPortal";
DialogClose.displayName = "DialogClose";
DialogOverlay.displayName = DialogPrimitive.Overlay.displayName;
DialogContent.displayName = DialogPrimitive.Content.displayName;
DialogHeader.displayName = "DialogHeader";
DialogFooter.displayName = "DialogFooter";
DialogTitle.displayName = DialogPrimitive.Title.displayName;
DialogDescription.displayName = DialogPrimitive.Description.displayName;

export {
	Dialog,
	DialogPortal,
	DialogOverlay,
	DialogTrigger,
	DialogClose,
	DialogContent,
	DialogHeader,
	DialogFooter,
	DialogTitle,
	DialogDescription,
};
