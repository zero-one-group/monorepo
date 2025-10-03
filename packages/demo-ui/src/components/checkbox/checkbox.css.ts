import { tv, type VariantProps } from "tailwind-variants/lite";

export const checkboxStyles = tv({
	slots: {
		base: [
			"peer size-4 shrink-0 rounded-sm border border-primary shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring",
			"disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground",
		],
		indicator: "flex items-center justify-center text-current",
		icon: "size-4",
	},
});

export type CheckboxVariants = VariantProps<typeof checkboxStyles>;
