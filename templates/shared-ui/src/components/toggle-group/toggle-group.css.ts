import { tv, type VariantProps } from "tailwind-variants/lite";

export const toggleGroupStyles = tv({
	base: "flex items-center justify-center gap-1",
});

export type ToggleGroupVariants = VariantProps<typeof toggleGroupStyles>;
