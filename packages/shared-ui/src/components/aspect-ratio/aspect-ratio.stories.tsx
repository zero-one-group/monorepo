import type { Meta, StoryObj } from "@storybook/react";
import { AspectRatio } from "./aspect-ratio";

const meta: Meta = {
	title: "Basic Components/AspectRatio",
	component: AspectRatio,
	argTypes: {
		ratio: {
			control: { type: "number" },
			description: "The desired aspect ratio",
			defaultValue: 16 / 9,
		},
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: (args) => (
		<div className="w-[450px]">
			<AspectRatio
				ratio={16 / 9}
				className="border border-white bg-muted"
				{...args}
			>
				<img
					src="https://images.unsplash.com/photo-1576075796033-848c2a5f3696?w=800&dpr=2&q=80"
					className="size-full rounded-md object-cover"
					alt="Pict by Alvaro Pinot"
				/>
			</AspectRatio>
		</div>
	),
};

export const Square: Story = {
	render: (args) => (
		<div className="w-[450px]">
			<AspectRatio ratio={1 / 1} className="bg-muted" {...args}>
				<img
					src="https://images.unsplash.com/photo-1588345921523-c2dcdb7f1dcd?w=800&dpr=2&q=80"
					alt="Pict by Drew Beamer"
					className="size-full rounded-md object-cover"
				/>
			</AspectRatio>
		</div>
	),
};

export const Portrait: Story = {
	render: (args) => (
		<div className="w-[300px]">
			<AspectRatio ratio={3 / 4} className="bg-muted" {...args}>
				<img
					src="https://images.unsplash.com/photo-1605656816944-971cd5c1407f?w=800&dpr=2&q=80"
					alt="Pict by Johannes Plenio"
					className="size-full rounded-md object-cover"
				/>
			</AspectRatio>
		</div>
	),
};
