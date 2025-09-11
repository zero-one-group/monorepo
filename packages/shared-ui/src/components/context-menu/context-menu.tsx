import * as Lucide from "lucide-react";
import { ContextMenu as ContextMenuPrimitive } from "radix-ui";
import * as React from "react";
import type { ContextMenuVariants } from "./context-menu.css";
import { contextMenuStyles } from "./context-menu.css";

const ContextMenu = ContextMenuPrimitive.Root;
const ContextMenuTrigger = ContextMenuPrimitive.Trigger;
const ContextMenuGroup = ContextMenuPrimitive.Group;
const ContextMenuPortal = ContextMenuPrimitive.Portal;
const ContextMenuSub = ContextMenuPrimitive.Sub;
const ContextMenuRadioGroup = ContextMenuPrimitive.RadioGroup;

interface ContextMenuSubTriggerProps
	extends React.ComponentPropsWithoutRef<
			typeof ContextMenuPrimitive.SubTrigger
		>,
		ContextMenuVariants {}

const ContextMenuSubTrigger = React.forwardRef<
	React.ComponentRef<typeof ContextMenuPrimitive.SubTrigger>,
	ContextMenuSubTriggerProps
>(({ className, inset, children, ...props }, ref) => {
	const styles = contextMenuStyles({ inset });
	return (
		<ContextMenuPrimitive.SubTrigger
			ref={ref}
			className={styles.subTrigger({ className })}
			{...props}
		>
			{children}
			<Lucide.ChevronRight className={styles.subTriggerIcon()} />
		</ContextMenuPrimitive.SubTrigger>
	);
});

const ContextMenuSubContent = React.forwardRef<
	React.ComponentRef<typeof ContextMenuPrimitive.SubContent>,
	React.ComponentPropsWithoutRef<typeof ContextMenuPrimitive.SubContent>
>(({ className, ...props }, ref) => {
	const styles = contextMenuStyles();
	return (
		<ContextMenuPrimitive.SubContent
			ref={ref}
			className={styles.subContent({ className })}
			{...props}
		/>
	);
});

const ContextMenuContent = React.forwardRef<
	React.ComponentRef<typeof ContextMenuPrimitive.Content>,
	React.ComponentPropsWithoutRef<typeof ContextMenuPrimitive.Content>
>(({ className, ...props }, ref) => {
	const styles = contextMenuStyles();
	return (
		<ContextMenuPrimitive.Portal>
			<ContextMenuPrimitive.Content
				ref={ref}
				className={styles.content({ className })}
				{...props}
			/>
		</ContextMenuPrimitive.Portal>
	);
});

interface ContextMenuItemProps
	extends React.ComponentPropsWithoutRef<typeof ContextMenuPrimitive.Item>,
		ContextMenuVariants {}

const ContextMenuItem = React.forwardRef<
	React.ComponentRef<typeof ContextMenuPrimitive.Item>,
	ContextMenuItemProps
>(({ className, inset, ...props }, ref) => {
	const styles = contextMenuStyles({ inset });
	return (
		<ContextMenuPrimitive.Item
			ref={ref}
			className={styles.item({ className })}
			{...props}
		/>
	);
});

const ContextMenuCheckboxItem = React.forwardRef<
	React.ComponentRef<typeof ContextMenuPrimitive.CheckboxItem>,
	React.ComponentPropsWithoutRef<typeof ContextMenuPrimitive.CheckboxItem>
>(({ className, children, checked, ...props }, ref) => {
	const styles = contextMenuStyles();
	return (
		<ContextMenuPrimitive.CheckboxItem
			ref={ref}
			className={styles.checkboxItem({ className })}
			checked={checked}
			{...props}
		>
			<ContextMenuPrimitive.ItemIndicator className={styles.itemIndicator()}>
				<Lucide.Check className={styles.itemIndicatorIcon()} />
			</ContextMenuPrimitive.ItemIndicator>
			{children}
		</ContextMenuPrimitive.CheckboxItem>
	);
});

const ContextMenuRadioItem = React.forwardRef<
	React.ComponentRef<typeof ContextMenuPrimitive.RadioItem>,
	React.ComponentPropsWithoutRef<typeof ContextMenuPrimitive.RadioItem>
>(({ className, children, ...props }, ref) => {
	const styles = contextMenuStyles();
	return (
		<ContextMenuPrimitive.RadioItem
			ref={ref}
			className={styles.radioItem({ className })}
			{...props}
		>
			<ContextMenuPrimitive.ItemIndicator className={styles.itemIndicator()}>
				<Lucide.Circle className={styles.radioItemIcon()} />
			</ContextMenuPrimitive.ItemIndicator>
			{children}
		</ContextMenuPrimitive.RadioItem>
	);
});

interface ContextMenuLabelProps
	extends React.ComponentPropsWithoutRef<typeof ContextMenuPrimitive.Label>,
		ContextMenuVariants {}

const ContextMenuLabel = React.forwardRef<
	React.ComponentRef<typeof ContextMenuPrimitive.Label>,
	ContextMenuLabelProps
>(({ className, inset, ...props }, ref) => {
	const styles = contextMenuStyles({ inset });
	return (
		<ContextMenuPrimitive.Label
			ref={ref}
			className={styles.label({ className })}
			{...props}
		/>
	);
});

const ContextMenuSeparator = React.forwardRef<
	React.ComponentRef<typeof ContextMenuPrimitive.Separator>,
	React.ComponentPropsWithoutRef<typeof ContextMenuPrimitive.Separator>
>(({ className, ...props }, ref) => {
	const styles = contextMenuStyles();
	return (
		<ContextMenuPrimitive.Separator
			ref={ref}
			className={styles.separator({ className })}
			{...props}
		/>
	);
});

const ContextMenuShortcut = ({
	className,
	...props
}: React.HTMLAttributes<HTMLSpanElement>) => {
	const styles = contextMenuStyles();
	return <span className={styles.shortcut({ className })} {...props} />;
};

export {
	ContextMenu,
	ContextMenuTrigger,
	ContextMenuContent,
	ContextMenuItem,
	ContextMenuCheckboxItem,
	ContextMenuRadioItem,
	ContextMenuLabel,
	ContextMenuSeparator,
	ContextMenuShortcut,
	ContextMenuGroup,
	ContextMenuPortal,
	ContextMenuSub,
	ContextMenuSubContent,
	ContextMenuSubTrigger,
	ContextMenuRadioGroup,
};
