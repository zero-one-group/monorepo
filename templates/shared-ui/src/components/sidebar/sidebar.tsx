import * as Lucide from "lucide-react";
import * as React from "react";
import { clx } from "../../utils";
import { Button } from "../button/button";
import { Sheet, SheetContent } from "../sheet/sheet";
import { useSidebar } from "./sidebar-provider";
import { sidebarStyles } from "./sidebar.css";

const SIDEBAR_WIDTH_MOBILE = "18rem";

const Sidebar = React.forwardRef<
	HTMLDivElement,
	React.ComponentProps<"div"> & {
		side?: "left" | "right";
		variant?: "sidebar" | "floating" | "inset";
		collapsible?: "offcanvas" | "icon" | "none";
	}
>(
	(
		{
			side = "left",
			variant = "sidebar",
			collapsible = "offcanvas",
			className,
			children,
			...props
		},
		ref,
	) => {
		const { isMobile, state, openMobile, setOpenMobile } = useSidebar();
		// const styles = sidebarStyles()

		if (collapsible === "none") {
			return (
				<div
					className={clx(
						"flex h-full w-(--sidebar-width) flex-col bg-sidebar text-sidebar-foreground",
						className,
					)}
					ref={ref}
					{...props}
				>
					{children}
				</div>
			);
		}

		if (isMobile) {
			return (
				<Sheet open={openMobile} onOpenChange={setOpenMobile} {...props}>
					<SheetContent
						data-sidebar="sidebar"
						data-mobile="true"
						className="w-(--sidebar-width) bg-sidebar p-0 text-sidebar-foreground [&>button]:hidden"
						style={
							{
								"--sidebar-width": SIDEBAR_WIDTH_MOBILE,
							} as React.CSSProperties
						}
						side={side}
					>
						<div className="flex size-full flex-col">{children}</div>
					</SheetContent>
				</Sheet>
			);
		}

		return (
			<div
				ref={ref}
				className="group peer hidden text-sidebar-foreground md:block"
				data-state={state}
				data-collapsible={state === "collapsed" ? collapsible : ""}
				data-variant={variant}
				data-side={side}
			>
				{/* This is what handles the sidebar gap on desktop */}
				<div
					className={clx(
						"relative h-svh w-(--sidebar-width) bg-transparent transition-[width] duration-200 ease-linear",
						"group-data-[collapsible=offcanvas]:w-0",
						"group-data-[side=right]:rotate-180",
						variant === "floating" || variant === "inset"
							? "group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)_+_theme(spacing.4))]"
							: "group-data-[collapsible=icon]:w-(--sidebar-width-icon)",
					)}
				/>
				<div
					className={clx(
						"fixed inset-y-0 z-10 hidden h-svh w-(--sidebar-width) transition-[left,right,width] duration-200 ease-linear md:flex",
						side === "left"
							? "left-0 group-data-[collapsible=offcanvas]:left-[calc(var(--sidebar-width)*-1)]"
							: "right-0 group-data-[collapsible=offcanvas]:right-[calc(var(--sidebar-width)*-1)]",
						// Adjust the padding for floating and inset variants.
						variant === "floating" || variant === "inset"
							? "p-2 group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)_+_theme(spacing.4)_+2px)]"
							: "group-data-[collapsible=icon]:w-(--sidebar-width-icon) group-data-[side=left]:border-r group-data-[side=right]:border-l",
						className,
					)}
					{...props}
				>
					<div
						data-sidebar="sidebar"
						className="flex size-full flex-col bg-sidebar group-data-[variant=floating]:rounded-lg group-data-[variant=floating]:border group-data-[variant=floating]:border-sidebar-border group-data-[variant=floating]:shadow"
					>
						{children}
					</div>
				</div>
			</div>
		);
	},
);

const SidebarTrigger = React.forwardRef<
	React.ComponentRef<typeof Button>,
	React.ComponentProps<typeof Button>
>(({ className, onClick, ...props }, ref) => {
	const { toggleSidebar, state } = useSidebar();
	const styles = sidebarStyles();

	const handleClick = (
		event: React.MouseEvent<HTMLButtonElement, MouseEvent>,
	) => {
		onClick?.(event);
		toggleSidebar();
	};

	return (
		<Button
			ref={ref}
			data-sidebar="trigger"
			variant="ghost"
			size="icon"
			className={styles.sidebarTriggerButton({ className })}
			onClick={handleClick}
			{...props}
		>
			{state === "expanded" ? (
				<Lucide.PanelLeftClose />
			) : (
				<Lucide.PanelLeftOpen />
			)}
			<span className="sr-only">Toggle Sidebar</span>
		</Button>
	);
});

Sidebar.displayName = "Sidebar";
SidebarTrigger.displayName = "SidebarTrigger";

export { Sidebar, SidebarTrigger };
