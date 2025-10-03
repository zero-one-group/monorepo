import { Command as CommandPrimitive } from "cmdk";
import * as Lucide from "lucide-react";
import type { Dialog as DialogPrimitive } from "radix-ui";
import * as React from "react";
import { Dialog, DialogContent } from "../dialog/dialog";
import { ScrollArea } from "../scroll-area/scroll-area";
import { commandStyles } from "./command.css";

const Command = React.forwardRef<
	React.ComponentRef<typeof CommandPrimitive>,
	React.ComponentPropsWithoutRef<typeof CommandPrimitive>
>(({ className, ...props }, ref) => {
	const styles = commandStyles();
	return (
		<CommandPrimitive
			ref={ref}
			className={styles.root({ className })}
			{...props}
		/>
	);
});

const CommandDialog = ({ children, ...props }: DialogPrimitive.DialogProps) => {
	const styles = commandStyles();
	return (
		<Dialog {...props}>
			<DialogContent className="overflow-hidden p-0">
				<Command className={styles.dialog()}>{children}</Command>
			</DialogContent>
		</Dialog>
	);
};

const CommandInput = React.forwardRef<
	React.ComponentRef<typeof CommandPrimitive.Input>,
	React.ComponentPropsWithoutRef<typeof CommandPrimitive.Input>
>(({ className, ...props }, ref) => {
	const styles = commandStyles();
	return (
		<div className={styles.inputWrapper()} cmdk-input-wrapper="">
			<Lucide.Search className={styles.searchIcon()} strokeWidth={2} />
			<CommandPrimitive.Input
				ref={ref}
				className={styles.input({ className })}
				{...props}
			/>
		</div>
	);
});

const CommandList = React.forwardRef<
	React.ComponentRef<typeof CommandPrimitive.List>,
	React.ComponentPropsWithoutRef<typeof CommandPrimitive.List>
>(({ className, ...props }, ref) => {
	const styles = commandStyles();
	return (
		<ScrollArea className={styles.list({ className })}>
			<CommandPrimitive.List
				className={styles.listInner()}
				ref={ref}
				{...props}
			/>
		</ScrollArea>
	);
});

const CommandEmpty = React.forwardRef<
	React.ComponentRef<typeof CommandPrimitive.Empty>,
	React.ComponentPropsWithoutRef<typeof CommandPrimitive.Empty>
>((props, ref) => {
	const styles = commandStyles();
	return (
		<CommandPrimitive.Empty ref={ref} className={styles.empty()} {...props} />
	);
});

const CommandGroup = React.forwardRef<
	React.ComponentRef<typeof CommandPrimitive.Group>,
	React.ComponentPropsWithoutRef<typeof CommandPrimitive.Group>
>(({ className, ...props }, ref) => {
	const styles = commandStyles();
	return (
		<CommandPrimitive.Group
			ref={ref}
			className={styles.group({ className })}
			{...props}
		/>
	);
});

const CommandSeparator = React.forwardRef<
	React.ComponentRef<typeof CommandPrimitive.Separator>,
	React.ComponentPropsWithoutRef<typeof CommandPrimitive.Separator>
>(({ className, ...props }, ref) => {
	const styles = commandStyles();
	return (
		<CommandPrimitive.Separator
			ref={ref}
			className={styles.separator({ className })}
			{...props}
		/>
	);
});

const CommandItem = React.forwardRef<
	React.ComponentRef<typeof CommandPrimitive.Item>,
	React.ComponentPropsWithoutRef<typeof CommandPrimitive.Item>
>(({ className, ...props }, ref) => {
	const styles = commandStyles();
	return (
		<CommandPrimitive.Item
			ref={ref}
			className={styles.item({ className })}
			{...props}
		/>
	);
});

const CommandShortcut = ({
	className,
	...props
}: React.HTMLAttributes<HTMLSpanElement>) => {
	const styles = commandStyles();
	return <span className={styles.shortcut({ className })} {...props} />;
};

Command.displayName = CommandPrimitive.displayName;
CommandInput.displayName = CommandPrimitive.Input.displayName;
CommandList.displayName = CommandPrimitive.List.displayName;
CommandEmpty.displayName = CommandPrimitive.Empty.displayName;
CommandGroup.displayName = CommandPrimitive.Group.displayName;
CommandSeparator.displayName = CommandPrimitive.Separator.displayName;
CommandItem.displayName = CommandPrimitive.Item.displayName;
CommandShortcut.displayName = "CommandShortcut";

export {
	Command,
	CommandDialog,
	CommandInput,
	CommandList,
	CommandEmpty,
	CommandGroup,
	CommandItem,
	CommandShortcut,
	CommandSeparator,
};
