import { tv, type VariantProps } from "tailwind-variants/lite";

export const sliderStyles = tv({
	slots: {
		base: "relative flex w-full touch-none select-none items-center",
		track:
			"relative h-1.5 w-full grow overflow-hidden rounded-full bg-primary/20",
		range: "absolute h-full bg-primary",
		thumb:
			"block size-4 rounded-full border border-primary/50 bg-background shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50",
	},
});

export type SliderVariants = VariantProps<typeof sliderStyles>;
