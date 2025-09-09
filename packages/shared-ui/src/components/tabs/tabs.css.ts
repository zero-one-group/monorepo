import { tv, type VariantProps } from "tailwind-variants";

export const tabsStyles = tv({
	slots: {
		list: "inline-flex h-9 items-center justify-center rounded-lg bg-muted p-1 text-muted-foreground",
		trigger: [
			"inline-flex cursor-pointer items-center justify-center whitespace-nowrap rounded-md px-3 py-1 font-medium text-sm",
			"ring-offset-background transition-all focus:ring-0 focus-visible:outline-none focus-visible:ring-1",
			"focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
			"data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow",
		],
		content: [
			"mt-2 ring-offset-background focus:ring-0 focus-visible:outline-none",
			"focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-2",
		],
	},
});

export type TabsVariants = VariantProps<typeof tabsStyles>;
