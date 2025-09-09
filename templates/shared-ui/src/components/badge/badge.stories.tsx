import type { Meta, StoryObj } from "@storybook/react";
import { Badge, type BadgeProps } from "./badge";
import type { BadgeVariants } from "./badge.css";

const variantOptions: NonNullable<BadgeVariants["variant"]>[] = [
	"default",
	"secondary",
	"success",
	"info",
	"warning",
	"destructive",
	"outline",
	"ghost",
];

const sizeOptions: NonNullable<BadgeVariants["size"]>[] = [
	"sm",
	"default",
	"lg",
];
const roundedOptions: NonNullable<BadgeVariants["rounded"]>[] = [
	"default",
	"full",
];

const meta: Meta<BadgeProps> = {
	title: "Basic Components/Badge",
	component: Badge,
	argTypes: {
		children: {
			control: "text",
			table: {
				type: { summary: "ReactNode | string" },
			},
		},
		variant: {
			control: { type: "radio" },
			options: variantOptions,
			table: {
				defaultValue: { summary: "default" },
				type: { summary: 'BadgeVariants["variant"]' },
			},
		},
		size: {
			control: { type: "inline-radio" },
			options: sizeOptions,
			table: {
				type: { summary: 'BadgeVariants["size"]' },
			},
		},
		rounded: {
			control: { type: "inline-radio" },
			options: roundedOptions,
		},
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		children: "Badge",
	},
};

export const VariantShowcase: Story = {
	parameters: {
		controls: { exclude: ["variant", "children"] },
	},
	render: (args) => (
		<div className="flex flex-wrap items-center gap-4">
			<Badge {...args}>Default</Badge>
			<Badge {...args} variant="secondary">
				Secondary
			</Badge>
			<Badge {...args} variant="success">
				Success
			</Badge>
			<Badge {...args} variant="info">
				Info
			</Badge>
			<Badge {...args} variant="warning">
				Warning
			</Badge>
			<Badge {...args} variant="destructive">
				Destructive
			</Badge>
			<Badge {...args} variant="outline">
				Outline
			</Badge>
			<Badge {...args} variant="ghost">
				Ghost
			</Badge>
		</div>
	),
};

export const SizeShowcase: Story = {
	parameters: {
		controls: { exclude: ["size", "children"] },
	},
	render: (args) => (
		<div className="flex flex-wrap items-end gap-4">
			<Badge {...args} size="sm">
				Small
			</Badge>
			<Badge {...args}>Default</Badge>
			<Badge {...args} size="lg">
				Large
			</Badge>
		</div>
	),
};

export const RoundedShowcase: Story = {
	parameters: {
		controls: { exclude: ["rounded", "children"] },
	},
	render: (args) => (
		<div className="flex flex-wrap items-center gap-4">
			<Badge {...args} rounded="default">
				Default Rounded
			</Badge>
			<Badge {...args} rounded="full">
				Full Rounded
			</Badge>
		</div>
	),
};
