import type { Meta, StoryObj } from "@storybook/react";
import {
	ContextMenu,
	ContextMenuCheckboxItem,
	ContextMenuContent,
	ContextMenuItem,
	ContextMenuLabel,
	ContextMenuRadioGroup,
	ContextMenuRadioItem,
	ContextMenuSeparator,
	ContextMenuShortcut,
	ContextMenuSub,
	ContextMenuSubContent,
	ContextMenuSubTrigger,
	ContextMenuTrigger,
} from "./context-menu";

const meta: Meta<typeof ContextMenu> = {
	title: "Basic Components/ContextMenu",
	component: ContextMenu,
	parameters: {
		layout: "centered",
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<ContextMenu>
			<ContextMenuTrigger className="flex h-[250px] w-[400px] items-center justify-center rounded-md border border-dashed text-sm">
				Right click here
			</ContextMenuTrigger>
			<ContextMenuContent className="w-64">
				<ContextMenuItem inset>
					Back
					<ContextMenuShortcut>⌘[</ContextMenuShortcut>
				</ContextMenuItem>
				<ContextMenuItem inset disabled>
					Forward
					<ContextMenuShortcut>⌘]</ContextMenuShortcut>
				</ContextMenuItem>
				<ContextMenuItem inset>
					Reload
					<ContextMenuShortcut>⌘R</ContextMenuShortcut>
				</ContextMenuItem>
				<ContextMenuSub>
					<ContextMenuSubTrigger inset>More Tools</ContextMenuSubTrigger>
					<ContextMenuSubContent className="w-48">
						<ContextMenuItem>
							Save Page As...
							<ContextMenuShortcut>⇧⌘S</ContextMenuShortcut>
						</ContextMenuItem>
						<ContextMenuItem>Create Shortcut...</ContextMenuItem>
						<ContextMenuItem>Name Window...</ContextMenuItem>
						<ContextMenuSeparator />
						<ContextMenuItem>Developer Tools</ContextMenuItem>
					</ContextMenuSubContent>
				</ContextMenuSub>
				<ContextMenuSeparator />
				<ContextMenuCheckboxItem checked>
					Show Bookmarks Bar
				</ContextMenuCheckboxItem>
				<ContextMenuCheckboxItem>Show Full URLs</ContextMenuCheckboxItem>
				<ContextMenuSeparator />
				<ContextMenuRadioGroup value="pedro">
					<ContextMenuLabel inset>People</ContextMenuLabel>
					<ContextMenuSeparator />
					<ContextMenuRadioItem value="pedro">
						Aris Ripandi
					</ContextMenuRadioItem>
					<ContextMenuRadioItem value="colm">Colm Tuite</ContextMenuRadioItem>
				</ContextMenuRadioGroup>
			</ContextMenuContent>
		</ContextMenu>
	),
};
