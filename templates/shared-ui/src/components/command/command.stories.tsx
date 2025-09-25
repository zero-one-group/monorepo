import type { Meta, StoryObj } from "@storybook/react";
import * as Lucide from "lucide-react";
import * as React from "react";
import {
	Command,
	CommandDialog,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
	CommandSeparator,
	CommandShortcut,
} from "./command";

const meta: Meta<typeof Command> = {
	title: "Basic Components/Command",
	component: Command,
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<Command className="rounded-lg border shadow-md">
			<CommandInput placeholder="Type a command or search..." />
			<CommandList className="min-w-[420px]">
				<CommandEmpty>No results found.</CommandEmpty>
				<CommandGroup heading="Suggestions">
					<CommandItem>
						<Lucide.Calendar className="mr-2 size-4" strokeWidth={2} />
						<span>Calendar</span>
					</CommandItem>
					<CommandItem>
						<Lucide.Smile className="mr-2 size-4" strokeWidth={2} />
						<span>Search Emoji</span>
					</CommandItem>
					<CommandItem>
						<Lucide.Calculator className="mr-2 size-4" strokeWidth={2} />
						<span>Calculator</span>
					</CommandItem>
				</CommandGroup>
				<CommandSeparator />
				<CommandGroup heading="Settings">
					<CommandItem>
						<Lucide.User className="mr-2 size-4" strokeWidth={2} />
						<span>Profile</span>
						<CommandShortcut>⌘P</CommandShortcut>
					</CommandItem>
					<CommandItem>
						<Lucide.CreditCard className="mr-2 size-4" strokeWidth={2} />
						<span>Billing</span>
						<CommandShortcut>⌘B</CommandShortcut>
					</CommandItem>
					<CommandItem>
						<Lucide.Settings className="mr-2 size-4" strokeWidth={2} />
						<span>Settings</span>
						<CommandShortcut>⌘S</CommandShortcut>
					</CommandItem>
				</CommandGroup>
			</CommandList>
		</Command>
	),
};

export const WithTrigger: Story = {
	render: () => {
		const [open, setOpen] = React.useState(false);

		React.useEffect(() => {
			const down = (e: KeyboardEvent) => {
				if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
					e.preventDefault();
					setOpen((open) => !open);
				}
			};
			document.addEventListener("keydown", down);
			return () => document.removeEventListener("keydown", down);
		}, []);

		return (
			<>
				<p className="text-muted-foreground text-sm">
					Press{" "}
					<kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded-sm border bg-muted px-1.5 font-medium font-mono text-[10px] text-muted-foreground opacity-100">
						<span className="text-xs">⌘</span>K
					</kbd>
				</p>
				<CommandDialog open={open} onOpenChange={setOpen}>
					<CommandInput placeholder="Type a command or search..." />
					<CommandList>
						<CommandEmpty>No results found.</CommandEmpty>
						<CommandGroup heading="Suggestions">
							<CommandItem>
								<Lucide.Calendar className="mr-2 size-4" strokeWidth={2} />
								<span>Calendar</span>
							</CommandItem>
							<CommandItem>
								<Lucide.Smile className="mr-2 size-4" strokeWidth={2} />
								<span>Search Emoji</span>
							</CommandItem>
							<CommandItem>
								<Lucide.Calculator className="mr-2 size-4" strokeWidth={2} />
								<span>Calculator</span>
							</CommandItem>
						</CommandGroup>
						<CommandSeparator />
						<CommandGroup heading="Settings">
							<CommandItem>
								<Lucide.User className="mr-2 size-4" strokeWidth={2} />
								<span>Profile</span>
								<CommandShortcut>⌘P</CommandShortcut>
							</CommandItem>
							<CommandItem>
								<Lucide.CreditCard className="mr-2 size-4" strokeWidth={2} />
								<span>Billing</span>
								<CommandShortcut>⌘B</CommandShortcut>
							</CommandItem>
							<CommandItem>
								<Lucide.Settings className="mr-2 size-4" strokeWidth={2} />
								<span>Settings</span>
								<CommandShortcut>⌘S</CommandShortcut>
							</CommandItem>
						</CommandGroup>
					</CommandList>
				</CommandDialog>
			</>
		);
	},
};
