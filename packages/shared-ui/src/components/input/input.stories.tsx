import type { Meta, StoryObj } from "@storybook/react-vite";
import { Input } from "./input";

const meta: Meta<typeof Input> = {
	title: "Basic Components/Input",
	component: Input,
	args: {
		placeholder: "Enter text...",
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Password: Story = {
	args: {
		type: "password",
		placeholder: "Enter password...",
	},
};

export const WithCopyButton: Story = {
	args: {
		value: "Text to copy",
		showCopyButton: true,
	},
};

export const WithExternalCopyButton: Story = {
	args: {
		value: "Text to copy",
		showExternalCopyButton: true,
	},
};

export const Disabled: Story = {
	args: {
		disabled: true,
	},
};
