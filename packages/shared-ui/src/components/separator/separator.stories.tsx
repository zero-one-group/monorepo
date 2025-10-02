import type { Meta, StoryObj } from "@storybook/react-vite";
import { Separator } from "./separator";
import type { SeparatorVariants } from "./separator.css";

const orientationOptions: NonNullable<SeparatorVariants["orientation"]>[] = [
	"horizontal",
	"vertical",
];

const meta: Meta<typeof Separator> = {
	title: "Basic Components/Separator",
	component: Separator,
	argTypes: {
		orientation: {
			control: "radio",
			options: orientationOptions,
			description: "The orientation of the separator",
			defaultValue: "horizontal",
		},
		decorative: {
			control: "boolean",
			description: "Whether the separator is purely decorative",
			defaultValue: true,
		},
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<div>
			<div className="space-y-1">
				<h4 className="font-medium text-sm leading-none">Radix Primitives</h4>
				<p className="text-muted-foreground text-sm">
					An open-source UI component library.
				</p>
			</div>
			<Separator className="my-4" />
			<div className="flex h-5 items-center space-x-4 text-sm">
				<div>Blog</div>
				<Separator orientation="vertical" />
				<div>Docs</div>
				<Separator orientation="vertical" />
				<div>Source</div>
			</div>
		</div>
	),
};

export const SeparatorShowcase: Story = {
	render: () => (
		<div className="space-y-8">
			<div className="space-y-4">
				<h4 className="font-medium text-sm">Horizontal Separator</h4>
				<div className="space-y-1">
					<h4 className="font-medium text-sm">Section 1</h4>
					<p className="text-muted-foreground text-sm">
						Content above separator
					</p>
				</div>
				<Separator />
				<div className="space-y-1">
					<h4 className="font-medium text-sm">Section 2</h4>
					<p className="text-muted-foreground text-sm">
						Content below separator
					</p>
				</div>
			</div>

			<div className="space-y-4">
				<h4 className="font-medium text-sm">Vertical Separator</h4>
				<div className="flex h-5 items-center space-x-4 text-sm">
					<div>Item 1</div>
					<Separator orientation="vertical" />
					<div>Item 2</div>
					<Separator orientation="vertical" />
					<div>Item 3</div>
				</div>
			</div>
		</div>
	),
};
