import type { Meta, StoryObj } from "@storybook/react-vite";
import { Input } from "../input/input";
import { Label } from "./label";

const meta: Meta<typeof Label> = {
	title: "Basic Components/Label",
	component: Label,
	argTypes: {
		children: {
			control: "text",
			table: {
				type: { summary: "ReactNode | string" },
			},
		},
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		children: "Email address",
	},
};

export const WithInput: Story = {
	render: (_args) => (
		<div className="grid w-full max-w-sm gap-1.5">
			<Label htmlFor="email">Email address</Label>
			<Input type="email" id="email" placeholder="Enter your email" />
		</div>
	),
};

export const Required: Story = {
	render: (_args) => (
		<div className="grid w-full max-w-sm gap-1.5">
			<Label htmlFor="username">
				Username <span className="text-destructive">*</span>
			</Label>
			<Input type="text" id="username" placeholder="Enter username" required />
		</div>
	),
};

export const Hidden: Story = {
	render: (_args) => (
		<div className="grid w-full max-w-sm gap-1.5">
			<Label htmlFor="username" hidden>
				Username <span className="text-destructive">*</span>
			</Label>
			<Input type="text" id="username" placeholder="Enter username" required />
		</div>
	),
};
