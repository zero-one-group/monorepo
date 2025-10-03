import * as Lucide from "lucide-react";
import { Checkbox as CheckboxPrimitive } from "radix-ui";
import * as React from "react";
import { checkboxStyles } from "./checkbox.css";

const Checkbox = React.forwardRef<
	React.ComponentRef<typeof CheckboxPrimitive.Root>,
	React.ComponentPropsWithoutRef<typeof CheckboxPrimitive.Root>
>(({ className, ...props }, ref) => {
	const styles = checkboxStyles();
	return (
		<CheckboxPrimitive.Root
			ref={ref}
			className={styles.base({ className })}
			{...props}
		>
			<CheckboxPrimitive.Indicator className={styles.indicator()}>
				<Lucide.Check className={styles.icon()} strokeWidth={2} />
			</CheckboxPrimitive.Indicator>
		</CheckboxPrimitive.Root>
	);
});

Checkbox.displayName = CheckboxPrimitive.Root.displayName;

export { Checkbox };
