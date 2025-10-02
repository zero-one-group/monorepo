import { type VariantProps, tv } from "tailwind-variants/lite";

export const drawerStyles = tv({
	slots: {
		overlay: "fixed inset-0 z-50 bg-black/80",
		content:
			"fixed inset-x-0 bottom-0 z-50 mt-24 flex h-auto flex-col rounded-t-[10px] border bg-background",
		handle: "mx-auto mt-4 h-2 w-[100px] rounded-full bg-muted",
		header: "grid gap-1.5 p-4 text-center sm:text-left",
		footer: "mt-auto flex flex-col gap-2 p-4",
		title: "font-semibold text-lg leading-none tracking-tight",
		description: "text-muted-foreground text-sm",
	},
});

export type DrawerVariants = VariantProps<typeof drawerStyles>;
