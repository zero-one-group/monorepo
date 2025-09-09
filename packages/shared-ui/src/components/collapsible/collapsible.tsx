import { Collapsible as CollapsiblePrimitive } from "radix-ui";
import * as React from "react";
import { collapsibleStyles } from "./collapsible.css";

const Collapsible = React.forwardRef<
	HTMLDivElement,
	React.ComponentPropsWithoutRef<typeof CollapsiblePrimitive.Root>
>(({ className, ...props }, ref) => {
	const styles = collapsibleStyles();
	return (
		<CollapsiblePrimitive.Root
			ref={ref}
			className={styles.base({ className })}
			{...props}
		/>
	);
});

const CollapsibleTrigger = React.forwardRef<
	HTMLButtonElement,
	React.ComponentPropsWithoutRef<typeof CollapsiblePrimitive.CollapsibleTrigger>
>(({ className, ...props }, ref) => {
	const styles = collapsibleStyles();
	return (
		<CollapsiblePrimitive.CollapsibleTrigger
			ref={ref}
			className={styles.trigger({ className })}
			{...props}
		/>
	);
});

const CollapsibleContent = React.forwardRef<
	HTMLDivElement,
	React.ComponentPropsWithoutRef<typeof CollapsiblePrimitive.CollapsibleContent>
>(({ className, ...props }, ref) => {
	const styles = collapsibleStyles();
	return (
		<CollapsiblePrimitive.CollapsibleContent
			ref={ref}
			className={styles.content({ className })}
			{...props}
		/>
	);
});

Collapsible.displayName = "Collapsible";
CollapsibleTrigger.displayName = "CollapsibleTrigger";
CollapsibleContent.displayName = "CollapsibleContent";

export { Collapsible, CollapsibleTrigger, CollapsibleContent };
