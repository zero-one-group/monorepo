import { tv, type VariantProps } from "tailwind-variants/lite";

export const toastStyles = tv({
	slots: {
		toast:
			"group toast group-[.toaster]:border-border group-[.toaster]:bg-background group-[.toaster]:text-foreground group-[.toaster]:shadow-lg",
		description: "group-[.toast]:text-muted-foreground",
		actionButton:
			"group-[.toast]:bg-primary group-[.toast]:text-primary-foreground",
		cancelButton:
			"group-[.toast]:bg-muted group-[.toast]:text-muted-foreground",
	},
});

export type ToastVariants = VariantProps<typeof toastStyles>;
