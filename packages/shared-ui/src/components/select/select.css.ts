import { tv, type VariantProps } from "tailwind-variants";

export const selectStyles = tv({
	slots: {
		trigger: [
			"flex h-9 w-full items-center justify-between whitespace-nowrap rounded-md border border-input",
			"bg-transparent px-3 py-2 text-sm shadow-xs ring-offset-background placeholder:text-muted-foreground",
			"focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
			"[&>span]:line-clamp-1",
		],
		icon: "size-4 opacity-50",
		content: [
			"relative z-50 max-h-96 min-w-[8rem] overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md",
			// Entry animations
			"data-[state=open]:motion-safe:motion-opacity-in-0",
			"data-[state=open]:motion-safe:motion-scale-in-95",
			"data-[state=open]:motion-safe:motion-duration-200",
			// Exit animations
			"data-[state=closed]:motion-safe:motion-opacity-out-0",
			"data-[state=closed]:motion-safe:motion-scale-out-95",
			"data-[state=closed]:motion-safe:motion-duration-150",
			// Slide animations based on position
			"data-[side=top]:motion-safe:motion-translate-y-in-2",
			"data-[side=bottom]:motion-safe:motion-translate-y-in--2",
			"data-[side=left]:motion-safe:motion-translate-x-in-2",
			"data-[side=right]:motion-safe:motion-translate-x-in--2",
		],
		viewport: "p-1",
		viewportPopper:
			"h-[var(--radix-select-trigger-height)] w-full min-w-[var(--radix-select-trigger-width)]",
		scrollButton: "flex cursor-default items-center justify-center py-1",
		label: "px-2 py-1.5 font-semibold text-sm",
		item: [
			"relative flex w-full cursor-default select-none items-center rounded-xs py-1.5 pr-8 pl-2",
			"text-sm outline-none focus:bg-accent focus:text-accent-foreground",
			"data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
		],
		itemIndicator: "absolute right-2 flex size-3.5 items-center justify-center",
		itemIndicatorIcon: "size-4",
		separator: "-mx-1 my-1 h-px bg-muted",
	},
});

export type SelectVariants = VariantProps<typeof selectStyles>;
