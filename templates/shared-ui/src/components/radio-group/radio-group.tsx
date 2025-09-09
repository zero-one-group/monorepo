import * as Lucide from "lucide-react";
import { RadioGroup as RadioGroupPrimitive } from "radix-ui";
import * as React from "react";
import { radioGroupStyles } from "./radio-group.css";
import type { RadioGroupVariants } from "./radio-group.css";

export interface RadioGroupProps
	extends React.ComponentPropsWithoutRef<typeof RadioGroupPrimitive.Root>,
		RadioGroupVariants {}

const RadioGroup = React.forwardRef<
	React.ComponentRef<typeof RadioGroupPrimitive.Root>,
	RadioGroupProps
>(({ className, orientation, ...props }, ref) => {
	const styles = radioGroupStyles({ orientation });
	return (
		<RadioGroupPrimitive.Root
			className={styles.root({ className })}
			{...props}
			ref={ref}
		/>
	);
});

const RadioGroupItem = React.forwardRef<
	React.ComponentRef<typeof RadioGroupPrimitive.Item>,
	React.ComponentPropsWithoutRef<typeof RadioGroupPrimitive.Item>
>(({ className, ...props }, ref) => {
	const styles = radioGroupStyles();
	return (
		<RadioGroupPrimitive.Item
			ref={ref}
			className={styles.item({ className })}
			{...props}
		>
			<RadioGroupPrimitive.Indicator className={styles.indicator()}>
				<Lucide.Check className={styles.icon()} strokeWidth={2} />
			</RadioGroupPrimitive.Indicator>
		</RadioGroupPrimitive.Item>
	);
});

RadioGroup.displayName = RadioGroupPrimitive.Root.displayName;
RadioGroupItem.displayName = RadioGroupPrimitive.Item.displayName;

export { RadioGroup, RadioGroupItem };
