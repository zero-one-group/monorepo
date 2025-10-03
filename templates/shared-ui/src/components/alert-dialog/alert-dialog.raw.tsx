import { AlertDialog as AlertDialogPrimitive } from "radix-ui";
import * as React from "react";
import { buttonStyles } from "../button/button.css";
import { alertDialogStyles } from "./alert-dialog.css";

const AlertDialog = AlertDialogPrimitive.Root;
const AlertDialogTrigger = AlertDialogPrimitive.Trigger;
const AlertDialogPortal = AlertDialogPrimitive.Portal;

const AlertDialogOverlay = React.forwardRef<
	React.ComponentRef<typeof AlertDialogPrimitive.Overlay>,
	React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Overlay>
>(({ className, ...props }, ref) => {
	const styles = alertDialogStyles();
	return (
		<AlertDialogPrimitive.Overlay
			className={styles.overlay({ className })}
			{...props}
			ref={ref}
		/>
	);
});

const AlertDialogContent = React.forwardRef<
	React.ComponentRef<typeof AlertDialogPrimitive.Content>,
	React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Content>
>(({ className, ...props }, ref) => {
	const styles = alertDialogStyles();
	return (
		<AlertDialogPortal>
			<AlertDialogOverlay />
			<AlertDialogPrimitive.Content
				ref={ref}
				className={styles.content({ className })}
				{...props}
			/>
		</AlertDialogPortal>
	);
});

const AlertDialogHeader = ({
	className,
	...props
}: React.HTMLAttributes<HTMLDivElement>) => {
	const styles = alertDialogStyles();
	return <div className={styles.header({ className })} {...props} />;
};

const AlertDialogFooter = ({
	className,
	...props
}: React.HTMLAttributes<HTMLDivElement>) => {
	const styles = alertDialogStyles();
	return <div className={styles.footer({ className })} {...props} />;
};

const AlertDialogTitle = React.forwardRef<
	React.ComponentRef<typeof AlertDialogPrimitive.Title>,
	React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Title>
>(({ className, ...props }, ref) => {
	const styles = alertDialogStyles();
	return (
		<AlertDialogPrimitive.Title
			ref={ref}
			className={styles.title({ className })}
			{...props}
		/>
	);
});

const AlertDialogDescription = React.forwardRef<
	React.ComponentRef<typeof AlertDialogPrimitive.Description>,
	React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Description>
>(({ className, ...props }, ref) => {
	const styles = alertDialogStyles();
	return (
		<AlertDialogPrimitive.Description
			ref={ref}
			className={styles.description({ className })}
			{...props}
		/>
	);
});

const AlertDialogAction = React.forwardRef<
	React.ComponentRef<typeof AlertDialogPrimitive.Action>,
	React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Action>
>(({ className, ...props }, ref) => {
	const styles = buttonStyles({ variant: "default", size: "default" });
	return (
		<AlertDialogPrimitive.Action
			ref={ref}
			className={styles.base({ className })}
			{...props}
		/>
	);
});

const AlertDialogCancel = React.forwardRef<
	React.ComponentRef<typeof AlertDialogPrimitive.Cancel>,
	React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Cancel>
>(({ className, ...props }, ref) => {
	const styles = buttonStyles({ variant: "ghost", size: "default" });
	return (
		<AlertDialogPrimitive.Cancel
			ref={ref}
			className={styles.base({ className })}
			{...props}
		/>
	);
});

AlertDialogOverlay.displayName = AlertDialogPrimitive.Overlay.displayName;
AlertDialogContent.displayName = AlertDialogPrimitive.Content.displayName;
AlertDialogHeader.displayName = "AlertDialogHeader";
AlertDialogFooter.displayName = "AlertDialogFooter";
AlertDialogTitle.displayName = AlertDialogPrimitive.Title.displayName;
AlertDialogDescription.displayName =
	AlertDialogPrimitive.Description.displayName;
AlertDialogAction.displayName = AlertDialogPrimitive.Action.displayName;
AlertDialogCancel.displayName = AlertDialogPrimitive.Cancel.displayName;

export {
	AlertDialog,
	AlertDialogPortal,
	AlertDialogOverlay,
	AlertDialogTrigger,
	AlertDialogContent,
	AlertDialogHeader,
	AlertDialogFooter,
	AlertDialogTitle,
	AlertDialogDescription,
	AlertDialogAction,
	AlertDialogCancel,
};
