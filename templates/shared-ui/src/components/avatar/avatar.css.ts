import { type VariantProps, tv } from "tailwind-variants/lite";

export const avatarStyles = tv({
	slots: {
		base: "relative flex size-10 shrink-0 overflow-hidden rounded-full",
		image: "aspect-square size-full",
		fallback:
			"flex size-full items-center justify-center rounded-full bg-muted",
	},
});

export type AvatarVariants = VariantProps<typeof avatarStyles>;
