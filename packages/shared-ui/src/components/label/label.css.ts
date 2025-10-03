import { tv, type VariantProps } from "tailwind-variants/lite";

export const labelStyles = tv({
	base: "font-medium text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
});

export type LabelVariants = VariantProps<typeof labelStyles>;
