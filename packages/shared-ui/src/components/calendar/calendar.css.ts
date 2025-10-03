import { tv, type VariantProps } from "tailwind-variants/lite";
import { buttonStyles } from "../button/button.css";

export const calendarStyles = tv({
	slots: {
		root: "p-3",
		months: "relative flex w-full flex-col items-center",
		month: "space-y-1",
		month_caption: "flex h-7 items-center justify-between",
		caption_label: "font-medium text-sm tracking-tight",
		nav: "absolute top-0 right-0 flex items-center justify-between space-x-1.5",
		button_previous: [
			buttonStyles({ variant: "outline" }).base(),
			"relative size-7 items-center justify-center p-0 opacity-50 hover:opacity-100",
		],
		button_next: [
			buttonStyles({ variant: "outline" }).base(),
			"relative size-7 items-center justify-center p-0 opacity-50 hover:opacity-100",
		],
		month_grid: "mt-4 w-full border-collapse space-y-1",
		weekdays: "flex",
		weekday: "w-8 rounded-sm font-normal text-[0.8rem] text-muted-foreground",
		week: "mt-2 flex w-full",
		cell: [
			"relative p-0 text-center text-sm focus-within:relative focus-within:z-20",
			"[&:has([aria-selected])]:bg-accent",
			"[&:has([aria-selected].day-outside)]:bg-accent/50",
			"[&:has([aria-selected].day-range-end)]:rounded-r-sm",
		],
		cell_range: [
			"[&:has(>.day-range-end)]:rounded-r-sm",
			"[&:has(>.day-range-start)]:rounded-l-sm",
			"first:[&:has([aria-selected])]:rounded-l-sm",
			"last:[&:has([aria-selected])]:rounded-r-sm",
		],
		cell_single: "[&:has([aria-selected])]:rounded-sm",
		day_button: [
			buttonStyles({ variant: "ghost" }).base(),
			"size-8 rounded-sm p-0 font-normal aria-selected:opacity-100",
		],
		range_start: "day-range-start",
		range_end: "day-range-end",
		selected:
			"rounded-sm bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground",
		today: "rounded-sm bg-accent text-accent-foreground",
		outside:
			"day-outside text-muted-foreground aria-selected:bg-accent/50 aria-selected:text-muted-foreground",
		disabled: "text-muted-foreground opacity-50",
		range_middle:
			"aria-selected:bg-accent aria-selected:text-accent-foreground",
		hidden: "invisible",
		icon: "size-4",
	},
});

export type CalendarVariants = VariantProps<typeof calendarStyles>;
