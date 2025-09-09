import { zodResolver } from "@hookform/resolvers/zod";
import type { Meta, StoryObj } from "@storybook/react";
import consola from "consola";
import { addDays, format } from "date-fns";
import * as Lucide from "lucide-react";
import * as React from "react";
import type { DateRange } from "react-day-picker";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { clx } from "../../utils";
import { Button } from "../button/button";
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "../form/form";
import { Popover, PopoverContent, PopoverTrigger } from "../popover/popover";
import { Calendar, type CalendarProps } from "./calendar";

const meta: Meta<CalendarProps> = {
	title: "Basic Components/Calendar",
	component: Calendar,
	parameters: {
		controls: {
			exclude: ["components", "classNames"],
		},
	},
};

export default meta;

export const Default: StoryObj = {
	args: {
		mode: "single",
	},
};

export const SimpleCalendar: StoryObj = {
	render: () => {
		const [date, setDate] = React.useState<Date | undefined>(new Date());

		return (
			<Calendar
				mode="single"
				selected={date}
				onSelect={setDate}
				className="w-full max-w-[250px] rounded-md border shadow"
			/>
		);
	},
};

const FormSchema = z.object({
	dob: z.date({
		required_error: "A date of birth is required.",
	}),
});

export const DatePickerWithForm: StoryObj = {
	render: () => {
		const form = useForm<z.infer<typeof FormSchema>>({
			resolver: zodResolver(FormSchema),
		});

		function onSubmit(data: z.infer<typeof FormSchema>) {
			consola.debug(data);
		}

		return (
			<Form {...form}>
				<form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
					<FormField
						control={form.control}
						name="dob"
						render={({ field }) => (
							<FormItem className="flex flex-col">
								<FormLabel>Date of birth</FormLabel>
								<Popover>
									<PopoverTrigger asChild>
										<FormControl>
											<Button
												variant="outline"
												className={clx(
													"w-[240px] pl-3 text-left font-normal",
													!field.value && "text-muted-foreground",
												)}
											>
												{field.value ? (
													format(field.value, "PPP")
												) : (
													<span>Pick a date</span>
												)}
												<Lucide.Calendar className="ml-auto size-4 opacity-50" />
											</Button>
										</FormControl>
									</PopoverTrigger>
									<PopoverContent className="w-auto p-0" align="start">
										<Calendar
											mode="single"
											selected={field.value}
											onSelect={field.onChange}
											disabled={(date) =>
												date > new Date() || date < new Date("1900-01-01")
											}
											initialFocus
										/>
									</PopoverContent>
								</Popover>
								<FormDescription>
									Your date of birth is used to calculate your age.
								</FormDescription>
								<FormMessage />
							</FormItem>
						)}
					/>
					<Button type="submit">Submit</Button>
				</form>
			</Form>
		);
	},
};

export const SelectionModes: StoryObj = {
	render: () => {
		const [singleDate, setSingleDate] = React.useState<Date | undefined>(
			new Date(),
		);
		const [rangeDate, setRangeDate] = React.useState<DateRange | undefined>({
			from: new Date(),
			to: addDays(new Date(), 7),
		});
		const [multiDates, setMultiDates] = React.useState<Date[]>([
			new Date(),
			addDays(new Date(), 2),
			addDays(new Date(), 5),
		]);

		return (
			<div className="flex flex-wrap items-start gap-4">
				<div>
					<h3 className="mb-4 font-medium text-sm">Single Date</h3>
					<Calendar
						mode="single"
						selected={singleDate}
						onSelect={setSingleDate}
						className="rounded-md border shadow"
					/>
				</div>
				<div>
					<h3 className="mb-4 font-medium text-sm">Date Range</h3>
					<Calendar
						mode="range"
						selected={rangeDate}
						onSelect={setRangeDate}
						className="rounded-md border shadow"
					/>
				</div>
				<div>
					<h3 className="mb-4 font-medium text-sm">Multiple Dates</h3>
					<Calendar
						mode="multiple"
						selected={multiDates}
						onSelect={(days: Date[] | undefined) => setMultiDates(days || [])}
						className="rounded-md border shadow"
					/>
				</div>
			</div>
		);
	},
};
