import { type VariantProps, tv } from "tailwind-variants";

export const dropdownMenuStyles = tv({
	slots: {
		trigger: [
			"flex cursor-default select-none items-center rounded-xs px-2 py-1.5 text-sm outline-none focus:bg-accent",
			"data-[state=open]:bg-accent",
		],
		triggerIcon: "ml-auto size-4",
		content: [
			"z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-sm",
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
		item: [
			"relative flex cursor-default select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none",
			"transition-colors focus:bg-accent focus:text-accent-foreground [&>svg]:size-4 [&>svg]:shrink-0",
			"data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
		],
		checkboxItem: [
			"relative flex cursor-default select-none items-center rounded-xs py-1.5 pr-2 pl-8 text-sm",
			"outline-none transition-colors focus:bg-accent focus:text-accent-foreground",
			"data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
		],
		checkboxItemIndicator:
			"absolute left-2 flex size-3.5 items-center justify-center",
		checkboxItemIcon: "size-4",
		radioItem: [
			"relative flex cursor-default select-none items-center rounded-xs py-1.5 pr-2 pl-8 text-sm",
			"outline-none transition-colors focus:bg-accent focus:text-accent-foreground",
			"data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
		],
		radioItemIndicator:
			"absolute left-2 flex size-3.5 items-center justify-center",
		radioItemIcon: "size-4 fill-current",
		label: "px-2 py-1.5 font-semibold text-sm",
		separator: "-mx-1 my-1 h-px bg-muted",
		shortcut: "ml-auto text-xs tracking-widest opacity-60",
		subTrigger: [
			"flex cursor-default select-none items-center rounded-xs px-2 py-1.5 text-sm outline-none focus:bg-accent",
			"data-[state=open]:bg-accent",
		],
		subContent: [
			"z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-sm",
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
	},
	variants: {
		inset: {
			true: {
				item: "pl-8",
				label: "pl-8",
				subTrigger: "pl-8",
			},
		},
	},
});

export type DropdownMenuVariants = VariantProps<typeof dropdownMenuStyles>;
