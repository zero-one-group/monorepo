import { Progress as ProgressPrimitive } from "radix-ui";
import * as React from "react";
import type { ProgressVariants } from "./progress.css";
import { progressStyles } from "./progress.css";

export interface ProgressProps
	extends React.ComponentPropsWithoutRef<typeof ProgressPrimitive.Root>,
		ProgressVariants {}

const Progress = React.forwardRef<
	React.ComponentRef<typeof ProgressPrimitive.Root>,
	ProgressProps
>(({ className, value, size, ...props }, ref) => {
	const styles = progressStyles({ size });
	return (
		<ProgressPrimitive.Root
			ref={ref}
			className={styles.base({ className })}
			{...props}
		>
			<ProgressPrimitive.Indicator
				className={styles.indicator()}
				style={{ transform: `translateX(-${100 - (value || 0)}%)` }}
			/>
		</ProgressPrimitive.Root>
	);
});

Progress.displayName = ProgressPrimitive.Root.displayName;

export { Progress };
