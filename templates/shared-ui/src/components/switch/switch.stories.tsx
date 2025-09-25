import type { Meta, StoryObj } from "@storybook/react";
import { Label } from "../label/label";
import { Switch } from "./switch";

const meta: Meta<typeof Switch> = {
	title: "Basic Components/Switch",
	component: Switch,
	argTypes: {
		checked: {
			control: "boolean",
			description: "The controlled checked state of the switch",
		},
		disabled: {
			control: "boolean",
			description:
				"When true, prevents the user from interacting with the switch",
		},
		defaultChecked: {
			control: "boolean",
			description: "The default checked state when initially rendered",
		},
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => <Switch />,
};

export const SwitchShowcase: Story = {
	render: () => (
		<div className="flex flex-col gap-4">
			<div className="flex items-center gap-2">
				<Switch id="airplane-mode" />
				<Label htmlFor="airplane-mode">Airplane Mode</Label>
			</div>

			<div className="flex items-center gap-2">
				<Switch id="disabled" disabled />
				<Label htmlFor="disabled">Disabled</Label>
			</div>

			<div className="flex items-center gap-2">
				<Switch id="default-checked" defaultChecked />
				<Label htmlFor="default-checked">Default Checked</Label>
			</div>
		</div>
	),
};
