import { tv, type VariantProps } from "tailwind-variants";

export const aspectRatioStyles = tv({
	base: "relative",
});

export type AspectRatioVariants = VariantProps<typeof aspectRatioStyles>;
