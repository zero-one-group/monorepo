import type { Meta, StoryObj } from "@storybook/react-vite";
import { Skeleton } from "./skeleton";

const meta: Meta<typeof Skeleton> = {
	title: "Basic Components/Skeleton",
	component: Skeleton,
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<div className="flex items-center space-x-4">
			<Skeleton className="size-12 rounded-full" />
			<div className="space-y-2">
				<Skeleton className="h-4 w-[250px]" />
				<Skeleton className="h-4 w-[200px]" />
			</div>
		</div>
	),
};

export const SkeletonShowcase: Story = {
	render: () => (
		<div className="space-y-8">
			<div className="space-y-4">
				<h4 className="font-medium text-sm">Profile Card Skeleton</h4>
				<div className="flex items-center space-x-4">
					<Skeleton className="size-12 rounded-full" />
					<div className="space-y-2">
						<Skeleton className="h-4 w-[250px]" />
						<Skeleton className="h-4 w-[200px]" />
					</div>
				</div>
			</div>

			<div className="space-y-4">
				<h4 className="font-medium text-sm">Text Block Skeleton</h4>
				<div className="space-y-3">
					<Skeleton className="h-4 w-full" />
					<Skeleton className="h-4 w-[90%]" />
					<Skeleton className="h-4 w-[80%]" />
				</div>
			</div>

			<div className="space-y-4">
				<h4 className="font-medium text-sm">Card Skeleton</h4>
				<div className="space-y-3">
					<Skeleton className="h-[200px] w-full rounded-xl" />
					<div className="space-y-2">
						<Skeleton className="h-4 w-[60%]" />
						<Skeleton className="h-4 w-[80%]" />
					</div>
				</div>
			</div>
		</div>
	),
};
