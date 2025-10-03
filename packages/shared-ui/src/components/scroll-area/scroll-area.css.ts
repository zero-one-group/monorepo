import { tv, type VariantProps } from "tailwind-variants/lite";

export const scrollAreaStyles = tv({
	slots: {
		root: "relative overflow-hidden",
		viewport: "size-full rounded-[inherit]",
		scrollbar: "flex touch-none select-none transition-colors",
		thumb: "relative flex-1 rounded-full bg-border",
	},
	variants: {
		orientation: {
			vertical: {
				scrollbar: "h-full w-2.5 border-l border-l-transparent p-[1px]",
			},
			horizontal: {
				scrollbar: "h-2.5 flex-col border-t border-t-transparent p-[1px]",
			},
		},
	},
	defaultVariants: {
		orientation: "vertical",
	},
});

export type ScrollAreaVariants = VariantProps<typeof scrollAreaStyles>;
