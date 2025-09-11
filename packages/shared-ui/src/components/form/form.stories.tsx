import type { Meta, StoryObj } from "@storybook/react";
import { useForm } from "react-hook-form";
import { Input } from "../input/input";
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "./form";

const meta: Meta<typeof Form> = {
	title: "Basic Components/Form",
	component: Form,
};

export default meta;
type Story = StoryObj<typeof meta>;

function FormDemo() {
	const form = useForm({
		defaultValues: {
			username: "",
		},
	});

	return (
		<Form {...form}>
			<form className="space-y-4">
				<FormField
					control={form.control}
					name="username"
					rules={{ required: "Username is required" }}
					render={({ field }) => (
						<FormItem>
							<FormLabel>Username</FormLabel>
							<FormControl>
								<Input placeholder="Enter username" {...field} />
							</FormControl>
							<FormDescription>
								This is your public display name.
							</FormDescription>
							<FormMessage />
						</FormItem>
					)}
				/>
			</form>
		</Form>
	);
}

export const Default: Story = {
	render: () => <FormDemo />,
};
