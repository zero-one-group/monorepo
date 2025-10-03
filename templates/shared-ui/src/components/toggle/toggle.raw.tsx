import { Toggle as TogglePrimitive } from "radix-ui";
import * as React from "react";
import type { ToggleVariants } from "./toggle.css";
import { toggleStyles } from "./toggle.css";

const Toggle = React.forwardRef<
	React.ComponentRef<typeof TogglePrimitive.Root>,
	React.ComponentPropsWithoutRef<typeof TogglePrimitive.Root> & ToggleVariants
>(({ className, variant, size, ...props }, ref) => (
	<TogglePrimitive.Root
		ref={ref}
		className={toggleStyles({ variant, size, className })}
		{...props}
	/>
));

Toggle.displayName = TogglePrimitive.Root.displayName;

export { Toggle };
