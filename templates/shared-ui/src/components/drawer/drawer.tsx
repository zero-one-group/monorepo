import * as React from "react";
import { Drawer as DrawerPrimitive } from "vaul";
import { drawerStyles } from "./drawer.css";

const Drawer = ({
	shouldScaleBackground = true,
	...props
}: React.ComponentProps<typeof DrawerPrimitive.Root>) => (
	<DrawerPrimitive.Root
		shouldScaleBackground={shouldScaleBackground}
		{...props}
	/>
);

type DrawerTriggerProps = React.ComponentProps<typeof DrawerPrimitive.Trigger>;
const DrawerTrigger: React.FC<DrawerTriggerProps> = DrawerPrimitive.Trigger;

type DrawerPortalProps = React.ComponentProps<typeof DrawerPrimitive.Portal>;
const DrawerPortal: React.FC<DrawerPortalProps> = DrawerPrimitive.Portal;

type DrawerCloseProps = React.ComponentProps<typeof DrawerPrimitive.Close>;
const DrawerClose: React.FC<DrawerCloseProps> = DrawerPrimitive.Close;

type DrawerOverlayProps = React.ComponentPropsWithoutRef<
	typeof DrawerPrimitive.Overlay
>;
type DrawerOverlayRef = React.ComponentRef<typeof DrawerPrimitive.Overlay>;

const DrawerOverlay: React.ForwardRefExoticComponent<
	DrawerOverlayProps & React.RefAttributes<DrawerOverlayRef>
> = React.forwardRef<DrawerOverlayRef, DrawerOverlayProps>(
	({ className, ...props }, ref) => {
		const styles = drawerStyles();
		return (
			<DrawerPrimitive.Overlay
				ref={ref}
				className={styles.overlay({ className })}
				{...props}
			/>
		);
	},
);

type DrawerContentProps = React.ComponentPropsWithoutRef<
	typeof DrawerPrimitive.Content
>;
type DrawerContentRef = React.ComponentRef<typeof DrawerPrimitive.Content>;

const DrawerContent: React.ForwardRefExoticComponent<
	DrawerContentProps & React.RefAttributes<DrawerContentRef>
> = React.forwardRef<DrawerContentRef, DrawerContentProps>(
	({ className, children, ...props }, ref) => {
		const styles = drawerStyles();
		return (
			<DrawerPortal>
				<DrawerOverlay />
				<DrawerPrimitive.Content
					ref={ref}
					className={styles.content({ className })}
					{...props}
				>
					<div className={styles.handle()} />
					{children}
				</DrawerPrimitive.Content>
			</DrawerPortal>
		);
	},
);

const DrawerHeader = ({
	className,
	...props
}: React.HTMLAttributes<HTMLDivElement>) => {
	const styles = drawerStyles();
	return <div className={styles.header({ className })} {...props} />;
};

const DrawerFooter = ({
	className,
	...props
}: React.HTMLAttributes<HTMLDivElement>) => {
	const styles = drawerStyles();
	return <div className={styles.footer({ className })} {...props} />;
};

type DrawerTitleProps = React.ComponentPropsWithoutRef<
	typeof DrawerPrimitive.Title
>;
type DrawerTitleRef = React.ComponentRef<typeof DrawerPrimitive.Title>;

const DrawerTitle: React.ForwardRefExoticComponent<
	DrawerTitleProps & React.RefAttributes<DrawerTitleRef>
> = React.forwardRef<DrawerTitleRef, DrawerTitleProps>(
	({ className, ...props }, ref) => {
		const styles = drawerStyles();
		return (
			<DrawerPrimitive.Title
				ref={ref}
				className={styles.title({ className })}
				{...props}
			/>
		);
	},
);

type DrawerDescriptionProps = React.ComponentPropsWithoutRef<
	typeof DrawerPrimitive.Description
>;
type DrawerDescriptionRef = React.ComponentRef<
	typeof DrawerPrimitive.Description
>;

const DrawerDescription: React.ForwardRefExoticComponent<
	DrawerDescriptionProps & React.RefAttributes<DrawerDescriptionRef>
> = React.forwardRef<DrawerDescriptionRef, DrawerDescriptionProps>(
	({ className, ...props }, ref) => {
		const styles = drawerStyles();
		return (
			<DrawerPrimitive.Description
				ref={ref}
				className={styles.description({ className })}
				{...props}
			/>
		);
	},
);

Drawer.displayName = "Drawer";
DrawerTrigger.displayName = "DrawerTrigger";
DrawerClose.displayName = "DrawerClose";
DrawerOverlay.displayName = "DrawerOverlay";
DrawerContent.displayName = "DrawerContent";
DrawerHeader.displayName = "DrawerHeader";
DrawerFooter.displayName = "DrawerFooter";
DrawerTitle.displayName = "DrawerTitle";
DrawerDescription.displayName = "DrawerDescription";

export {
	Drawer,
	DrawerPortal,
	DrawerOverlay,
	DrawerTrigger,
	DrawerClose,
	DrawerContent,
	DrawerHeader,
	DrawerFooter,
	DrawerTitle,
	DrawerDescription,
};
