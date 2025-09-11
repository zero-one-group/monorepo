import { tv, type VariantProps } from "tailwind-variants";

export const textStyles = tv({
	base: "transition-colors duration-200 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring",
	variants: {
		size: {
			xs: "text-xs",
			sm: "text-sm",
			md: "text-base",
			lg: "text-lg",
			xl: "text-xl",
			"2xl": "text-2xl",
		},
		weight: {
			light: "font-light",
			normal: "font-normal",
			medium: "font-medium",
			semibold: "font-semibold",
			bold: "font-bold",
		},
		align: {
			left: "text-left",
			center: "text-center",
			right: "text-right",
		},
		variant: {
			default: "text-foreground hover:text-foreground/90",
			muted: "text-muted-foreground hover:text-muted-foreground/90",
			primary: "text-primary hover:text-primary/90",
			secondary: "text-secondary-foreground hover:text-secondary-foreground/90",
			destructive: "text-destructive hover:text-destructive/90",
			success: "text-success hover:text-success/90",
			warning: "text-warning hover:text-warning/90",
		},
	},
	compoundVariants: [
		{
			size: ["lg", "xl", "2xl"],
			weight: ["semibold", "bold"],
			className: "leading-normal tracking-tight",
		},
		{
			size: ["xs", "sm"],
			weight: ["light", "normal"],
			className: "leading-tight tracking-wide",
		},
	],
	defaultVariants: {
		size: "md",
		weight: "normal",
		align: "left",
		variant: "default",
	},
});

export type TextVariants = VariantProps<typeof textStyles>;
