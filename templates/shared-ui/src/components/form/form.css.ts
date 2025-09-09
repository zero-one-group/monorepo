import { type VariantProps, tv } from "tailwind-variants";

export const formStyles = tv({
	slots: {
		item: "space-y-2",
		label:
			"font-medium text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
		labelError: "text-destructive",
		control: "",
		description: "text-[0.8rem] text-muted-foreground",
		message: "font-medium text-[0.8rem] text-destructive",
	},
	variants: {
		error: {
			true: {
				label: "text-destructive",
			},
		},
	},
});

export type FormVariants = VariantProps<typeof formStyles>;
