import { tv, type VariantProps } from "tailwind-variants";

export const badgeStyles = tv({
	slots: {
		base: "inline-flex items-center rounded-md border font-medium text-xs transition-colors",
	},
	variants: {
		variant: {
			default:
				"border-transparent bg-primary text-primary-foreground shadow-sm hover:bg-primary/80",
			secondary:
				"border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
			success:
				"border-transparent bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-300",
			destructive:
				"border-transparent bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/80",
			outline:
				"border-input bg-background text-foreground/80 hover:bg-accent/50 hover:text-accent-foreground",
			info: "border-transparent bg-sky-100 text-sky-700 dark:bg-sky-500/20 dark:text-sky-300",
			warning:
				"border-transparent bg-amber-100 text-amber-700 dark:bg-amber-500/20 dark:text-amber-300",
			ghost: "border-transparent text-muted-foreground hover:bg-muted/50",
		},
		size: {
			sm: "px-2 py-0.5 text-[10px]",
			default: "px-2.5 py-0.5 text-xs",
			lg: "px-3 py-1 text-sm",
		},
		rounded: {
			default: "rounded-md",
			full: "rounded-full",
		},
	},
	defaultVariants: {
		variant: "default",
		size: "default",
		rounded: "default",
	},
});

export type BadgeVariants = VariantProps<typeof badgeStyles>;
