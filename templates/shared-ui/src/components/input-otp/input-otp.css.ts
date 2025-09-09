import { type VariantProps, tv } from "tailwind-variants";

export const inputOtpStyles = tv({
	slots: {
		root: "flex items-center gap-2 has-[:disabled]:opacity-50",
		input: "disabled:cursor-not-allowed",
		group: "flex items-center",
		slot: "relative flex items-center justify-center border-input border-y border-r text-sm shadow-xs transition-all first:rounded-l-md first:border-l last:rounded-r-md",
		slotActive: "z-10 ring-1 ring-ring",
		caret:
			"pointer-events-none absolute inset-0 flex items-center justify-center",
		caretBlink: "motion-preset-blink motion-duration-1000 w-px bg-foreground",
		separator: "mx-1 [&>svg]:shrink-0",
	},
	variants: {
		size: {
			sm: {
				slot: "size-7 text-xs",
				caretBlink: "h-3",
				separator: "[&>svg]:size-3",
			},
			default: {
				slot: "size-9 text-sm",
				caretBlink: "h-4",
				separator: "[&>svg]:size-4",
			},
			lg: {
				slot: "size-11 text-base",
				caretBlink: "h-5",
				separator: "[&>svg]:size-5",
			},
		},
	},
	defaultVariants: {
		size: "default",
	},
});

export type InputOtpVariants = VariantProps<typeof inputOtpStyles>;
