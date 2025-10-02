import { type VariantProps, tv } from "tailwind-variants/lite";

export const accordionStyles = tv({
	slots: {
		item: "border border-b-0 px-4 first:rounded-t-md last:rounded-b-md last:border-b",
		headerWrapper: "flex",
		headerTrigger: [
			"flex flex-1 cursor-pointer items-center justify-between py-4 font-medium text-sm transition-all hover:underline",
			"[&[data-state=open]>svg]:rotate-180",
		],
		headerIcon:
			"size-4 shrink-0 text-muted-foreground transition-transform duration-200",
		contentWrapper: [
			"overflow-hidden text-sm",
			"data-[state=closed]:animate-accordion-up",
			"data-[state=open]:animate-accordion-down",
		],
		contentChildren: "pt-0 pb-4",
	},
});

export type AccordionVariants = VariantProps<typeof accordionStyles>;
