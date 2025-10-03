import type { Meta, StoryObj } from "@storybook/react-vite";
import { Heading, type HeadingProps } from "./heading";
import type { HeadingVariants } from "./heading.css";

const levelOptions: NonNullable<HeadingVariants["level"]>[] = [
	"h1",
	"h2",
	"h3",
	"h4",
	"h5",
	"h6",
];
const weightOptions: NonNullable<HeadingVariants["weight"]>[] = [
	"light",
	"normal",
	"medium",
	"semibold",
	"bold",
	"extrabold",
];
const alignOptions: NonNullable<HeadingVariants["align"]>[] = [
	"left",
	"center",
	"right",
];

const meta: Meta<HeadingProps> = {
	title: "Basic Components/Heading",
	component: Heading,
	parameters: {
		controls: {
			exclude: ["asChild"],
		},
	},
	argTypes: {
		children: {
			control: "text",
			description: "Heading content",
			table: {
				type: { summary: "ReactNode | string" },
			},
		},
		level: {
			control: { type: "select" },
			options: levelOptions,
			table: {
				type: { summary: 'HeadingVariants["level"]' },
			},
		},
		weight: {
			control: { type: "select" },
			options: weightOptions,
			table: {
				type: { summary: 'HeadingVariants["weight"]' },
			},
		},
		align: {
			control: { type: "radio" },
			options: alignOptions,
			table: {
				type: { summary: 'HeadingVariants["align"]' },
			},
		},
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		children: "The quick brown fox jumps over the lazy dog",
	},
};

export const LevelShowcase: Story = {
	parameters: {
		controls: { exclude: ["level"] },
	},
	render: () => (
		<div className="flex flex-col gap-4">
			<Heading level="h1">Heading Level 1</Heading>
			<Heading level="h2">Heading Level 2</Heading>
			<Heading level="h3">Heading Level 3</Heading>
			<Heading level="h4">Heading Level 4</Heading>
			<Heading level="h5">Heading Level 5</Heading>
			<Heading level="h6">Heading Level 6</Heading>
		</div>
	),
};

export const WeightShowcase: Story = {
	parameters: {
		controls: { exclude: ["weight"] },
	},
	render: () => (
		<div className="flex flex-col gap-4">
			<Heading weight="light">Light weight heading</Heading>
			<Heading weight="normal">Normal weight heading</Heading>
			<Heading weight="medium">Medium weight heading</Heading>
			<Heading weight="semibold">Semibold weight heading</Heading>
			<Heading weight="bold">Bold weight heading</Heading>
			<Heading weight="extrabold">Extrabold weight heading</Heading>
		</div>
	),
};

export const AlignmentShowcase: Story = {
	parameters: {
		controls: { exclude: ["align"] },
	},
	render: () => (
		<div className="flex flex-col gap-4">
			<Heading align="left">Left aligned heading</Heading>
			<Heading align="center">Center aligned heading</Heading>
			<Heading align="right">Right aligned heading</Heading>
		</div>
	),
};
