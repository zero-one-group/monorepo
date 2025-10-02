import { type VariantProps, tv } from "tailwind-variants/lite";

export const toggleStyles = tv({
	base: "inline-flex items-center justify-center gap-2 rounded-md font-medium text-sm transition-colors hover:bg-muted hover:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 data-[state=on]:bg-accent data-[state=on]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0",
	variants: {
		variant: {
			default: "bg-transparent",
			outline:
				"border border-input bg-transparent shadow-xs hover:bg-accent hover:text-accent-foreground",
		},
		size: {
			default: "h-9 min-w-9 px-2",
			sm: "h-8 min-w-8 px-1.5",
			lg: "h-10 min-w-10 px-2.5",
		},
	},
	compoundVariants: [],
	defaultVariants: {
		variant: "default",
		size: "default",
	},
});

export type ToggleVariants = VariantProps<typeof toggleStyles>;
