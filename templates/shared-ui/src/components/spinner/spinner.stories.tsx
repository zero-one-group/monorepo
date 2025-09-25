import type { Meta, StoryObj } from "@storybook/react";
import { LazyLoadingSpinner, LoadingSpinner } from "./spinner";

const meta: Meta<typeof LoadingSpinner> = {
	title: "Basic Components/Spinner",
	component: LoadingSpinner,
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => <LoadingSpinner />,
};

export const LazySpinner: Story = {
	render: () => <LazyLoadingSpinner delay={1000} />,
};

export const SpinnerShowcase: Story = {
	render: () => (
		<div className="flex items-center gap-4">
			<LoadingSpinner />
			<LazyLoadingSpinner delay={1000} />
			<LazyLoadingSpinner
				delay={1500}
				element={<LoadingSpinner />}
				fallback={<div>Loading...</div>}
			/>
		</div>
	),
};
