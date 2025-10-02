import { tv, type VariantProps } from "tailwind-variants/lite";

export const headingStyles = tv({
	base: "font-heading text-foreground tracking-tight transition-colors duration-200",
	variants: {
		level: {
			h1: "font-extrabold text-4xl lg:text-5xl",
			h2: "font-bold text-3xl lg:text-4xl",
			h3: "font-bold text-2xl lg:text-3xl",
			h4: "font-semibold text-xl lg:text-2xl",
			h5: "font-semibold text-lg lg:text-xl",
			h6: "font-semibold text-base lg:text-lg",
		},
		weight: {
			light: "font-light",
			normal: "font-normal",
			medium: "font-medium",
			semibold: "font-semibold",
			bold: "font-bold",
			extrabold: "font-extrabold",
		},
		align: {
			left: "text-left",
			center: "text-center",
			right: "text-right",
		},
	},
	defaultVariants: {
		level: "h2",
		align: "left",
	},
});

export type HeadingVariants = VariantProps<typeof headingStyles>;
