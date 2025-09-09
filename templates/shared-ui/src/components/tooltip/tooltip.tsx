import { Tooltip as TooltipPrimitive } from "radix-ui";
import * as React from "react";
import { tooltipStyles } from "./tooltip.css";

const TooltipProvider = TooltipPrimitive.Provider;
const Tooltip = TooltipPrimitive.Root;
const TooltipTrigger = TooltipPrimitive.Trigger;

const TooltipContent = React.forwardRef<
	React.ComponentRef<typeof TooltipPrimitive.Content>,
	React.ComponentPropsWithoutRef<typeof TooltipPrimitive.Content>
>(({ className, sideOffset = 4, ...props }, ref) => {
	const styles = tooltipStyles();
	return (
		<TooltipPrimitive.Portal>
			<TooltipPrimitive.Content
				ref={ref}
				sideOffset={sideOffset}
				className={styles.base({ className })}
				{...props}
			>
				{props.children}
				<TooltipPrimitive.Arrow className={styles.arrow()} />
			</TooltipPrimitive.Content>
		</TooltipPrimitive.Portal>
	);
});

TooltipContent.displayName = TooltipPrimitive.Content.displayName;

export { Tooltip, TooltipTrigger, TooltipContent, TooltipProvider };
