import { Separator as SeparatorPrimitive } from "radix-ui";
import * as React from "react";
import { type SeparatorVariants, separatorStyles } from "./separator.css";

const Separator = React.forwardRef<
	React.ComponentRef<typeof SeparatorPrimitive.Root>,
	React.ComponentPropsWithoutRef<typeof SeparatorPrimitive.Root> &
		SeparatorVariants
>(
	(
		{ className, orientation = "horizontal", decorative = true, ...props },
		ref,
	) => (
		<SeparatorPrimitive.Root
			ref={ref}
			decorative={decorative}
			orientation={orientation}
			className={separatorStyles({ orientation, className })}
			{...props}
		/>
	),
);

Separator.displayName = SeparatorPrimitive.Root.displayName;

export { Separator };
