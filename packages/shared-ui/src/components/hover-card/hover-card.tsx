import { HoverCard as HoverCardPrimitive } from "radix-ui";
import * as React from "react";
import type { HoverCardVariants } from "./hover-card.css";
import { hoverCardStyles } from "./hover-card.css";

const HoverCard = HoverCardPrimitive.Root;
const HoverCardTrigger = HoverCardPrimitive.Trigger;

interface HoverCardContentProps
	extends React.ComponentPropsWithoutRef<typeof HoverCardPrimitive.Content>,
		HoverCardVariants {}

const HoverCardContent = React.forwardRef<
	React.ComponentRef<typeof HoverCardPrimitive.Content>,
	HoverCardContentProps
>(({ className, align = "center", sideOffset = 4, ...props }, ref) => {
	const styles = hoverCardStyles();
	return (
		<HoverCardPrimitive.Content
			ref={ref}
			align={align}
			sideOffset={sideOffset}
			className={styles.content({ className })}
			{...props}
		/>
	);
});

HoverCard.displayName = HoverCardPrimitive.Root.displayName;
HoverCardTrigger.displayName = HoverCardPrimitive.Trigger.displayName;
HoverCardContent.displayName = HoverCardPrimitive.Content.displayName;

export { HoverCard, HoverCardTrigger, HoverCardContent };
