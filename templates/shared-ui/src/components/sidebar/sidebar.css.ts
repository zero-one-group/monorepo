import { type VariantProps, tv } from "tailwind-variants";

export const sidebarStyles = tv({
	slots: {
		sidebarProvider:
			"group/sidebar-wrapper flex min-h-svh w-full has-[[data-variant=inset]]:bg-sidebar",
		sidebarGroup: "relative flex w-full min-w-0 flex-col p-2",
		sidebarGroupLabel: [
			"flex h-8 shrink-0 items-center rounded-md px-2 font-medium text-sidebar-foreground/70 text-xs outline-none ring-sidebar-ring transition-[margin,opa] duration-200 ease-linear focus:ring-0 focus-visible:ring-1 [&>svg]:size-4 [&>svg]:shrink-0",
			"group-data-[collapsible=icon]:-mt-8 group-data-[collapsible=icon]:opacity-0",
		],
		sidebarGroupAction: [
			"absolute top-3.5 right-3 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground outline-none ring-sidebar-ring transition-transform hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus:ring-0 focus-visible:ring-1 [&>svg]:size-4 [&>svg]:shrink-0",
			// Increases the hit area of the button on mobile.
			"after:-inset-2 after:absolute after:md:hidden",
			"group-data-[collapsible=icon]:hidden",
		],
		sidebarGroupContent: "w-full text-sm",
		sidebarTriggerButton: "size-8",
		sidebarMenu: "flex w-full min-w-0 flex-col gap-1",
		sidebarMenuItem: "group/menu-item relative",
		sidebarMenuAction: [
			"absolute top-1.5 right-1 flex aspect-square w-5 items-center justify-center rounded-md p-0 text-sidebar-foreground outline-none ring-sidebar-ring transition-transform hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus:ring-0 focus-visible:ring-1 peer-hover/menu-button:text-sidebar-accent-foreground [&>svg]:size-4 [&>svg]:shrink-0",
			// Increases the hit area of the button on mobile.
			"after:-inset-2 after:absolute after:md:hidden",
			"peer-data-[size=sm]/menu-button:top-1",
			"peer-data-[size=default]/menu-button:top-1.5",
			"peer-data-[size=lg]/menu-button:top-2.5",
			"group-data-[collapsible=icon]:hidden",
		],
		sidebarMenuBadge: [
			"pointer-events-none absolute right-1 flex h-5 min-w-5 select-none items-center justify-center rounded-md px-1 font-medium text-sidebar-foreground text-xs tabular-nums",
			"peer-hover/menu-button:text-sidebar-accent-foreground peer-data-[active=true]/menu-button:text-sidebar-accent-foreground",
			"peer-data-[size=sm]/menu-button:top-1",
			"peer-data-[size=default]/menu-button:top-1.5",
			"peer-data-[size=lg]/menu-button:top-2.5",
			"group-data-[collapsible=icon]:hidden",
		],
		sidebarMenuSkeletonWrapper: "flex h-8 items-center gap-2 rounded-md px-2",
		sidebarMenuSkeletonIcon: "size-4 rounded-md",
		sidebarMenuSkeletonText: "h-4 max-w-(--skeleton-width) flex-1",
		sidebarMenuSub: [
			"mx-3.5 flex min-w-0 translate-x-px flex-col gap-1 border-sidebar-border border-l px-2.5 py-0.5",
			"group-data-[collapsible=icon]:hidden",
		],
	},
	variants: {
		showOnHover: {
			true: "group-focus-within/menu-item:opacity-100 group-hover/menu-item:opacity-100 data-[state=open]:opacity-100 peer-data-[active=true]/menu-button:text-sidebar-accent-foreground md:opacity-0",
			false: "",
		},
	},
	compoundVariants: [],
	defaultVariants: {
		showOnHover: false,
	},
});

export type SidebarVariants = VariantProps<typeof sidebarStyles>;

export const sidebarMenuButtonStyles = tv({
	base: "peer/menu-button group-data-[collapsible=icon]:!size-8 group-data-[collapsible=icon]:!p-2 flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left text-sm outline-none ring-sidebar-ring transition-[width,height,padding] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus:ring-0 focus-visible:ring-1 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 group-has-[[data-sidebar=menu-action]]/menu-item:pr-8 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0",
	variants: {
		variant: {
			default: "hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
			outline:
				"bg-background shadow-[0_0_0_1px_hsl(var(--sidebar-border))] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:shadow-[0_0_0_1px_hsl(var(--sidebar-accent))]",
		},
		size: {
			default: "h-8 text-sm",
			sm: "h-7 text-xs",
			lg: "group-data-[collapsible=icon]:!p-0 h-12 text-sm",
		},
	},
	compoundVariants: [],
	defaultVariants: {
		variant: "default",
		size: "default",
	},
});

export type SidebarMenuButtonVariants = VariantProps<
	typeof sidebarMenuButtonStyles
>;
