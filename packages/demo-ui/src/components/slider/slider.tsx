import { Slider as SliderPrimitive } from "radix-ui";
import * as React from "react";
import { sliderStyles } from "./slider.css";

const Slider = React.forwardRef<
	React.ComponentRef<typeof SliderPrimitive.Root>,
	React.ComponentPropsWithoutRef<typeof SliderPrimitive.Root>
>(({ className, ...props }, ref) => {
	const styles = sliderStyles();
	return (
		<SliderPrimitive.Root
			ref={ref}
			className={styles.base({ className })}
			{...props}
		>
			<SliderPrimitive.Track className={styles.track()}>
				<SliderPrimitive.Range className={styles.range()} />
			</SliderPrimitive.Track>
			<SliderPrimitive.Thumb className={styles.thumb()} />
		</SliderPrimitive.Root>
	);
});

Slider.displayName = SliderPrimitive.Root.displayName;

export { Slider };
