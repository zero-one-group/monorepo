import type { Meta, StoryObj } from "@storybook/react";
import { Separator } from "../separator/separator";
import { ScrollArea } from "./scroll-area";

const meta: Meta<typeof ScrollArea> = {
	title: "Basic Components/ScrollArea",
	component: ScrollArea,
};

export default meta;
type Story = StoryObj<typeof meta>;

const tags = Array.from({ length: 50 }).map(
	(_, i, a) => `v1.2.0-beta.${a.length - i}`,
);

export const Default: Story = {
	render: () => (
		<ScrollArea className="h-72 w-48 rounded-md border">
			<div className="p-4">
				<h4 className="mb-4 font-medium text-sm leading-none">Tags</h4>
				{tags.map((tag) => (
					<>
						<div key={tag} className="text-sm">
							{tag}
						</div>
						<Separator key={tag} className="my-2" />
					</>
				))}
			</div>
		</ScrollArea>
	),
};
