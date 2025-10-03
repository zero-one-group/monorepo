import { tv, type VariantProps } from "tailwind-variants/lite";

export const breadcrumbStyles = tv({
	slots: {
		base: "",
		list: "flex flex-wrap items-center gap-1.5 break-words text-muted-foreground text-sm sm:gap-2.5",
		item: "inline-flex items-center gap-1.5",
		link: "transition-colors hover:text-foreground",
		page: "font-normal text-foreground",
		separator: "[&>svg]:h-3.5 [&>svg]:w-3.5",
		ellipsis: "flex size-9 cursor-pointer items-center justify-center",
		ellipsisIcon: "size-4",
	},
});

export type BreadcrumbVariants = VariantProps<typeof breadcrumbStyles>;
