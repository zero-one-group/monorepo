import type { Meta, StoryObj } from "@storybook/react-vite";
import { Progress } from "./progress";
import type { ProgressVariants } from "./progress.css";

const sizeOptions: NonNullable<ProgressVariants["size"]>[] = [
	"default",
	"sm",
	"lg",
];

const meta: Meta<typeof Progress> = {
	title: "Basic Components/Progress",
	component: Progress,
	args: {
		value: 60,
	},
	argTypes: {
		value: {
			control: { type: "range", min: 0, max: 100 },
		},
		size: {
			control: "radio",
			options: sizeOptions,
			table: {
				defaultValue: { summary: "default" },
				type: { summary: 'ProgressVariants["size"]' },
			},
		},
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: (args) => (
		<div className="min-w-[300px]">
			<Progress {...args} />
		</div>
	),
};

export const Sizes: Story = {
	render: (args) => (
		<div className="min-w-[300px] space-y-4">
			<Progress size="sm" value={args.value} />
			<Progress size="default" value={args.value} />
			<Progress size="lg" value={args.value} />
		</div>
	),
};
