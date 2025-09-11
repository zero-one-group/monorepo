import { tv, type VariantProps } from "tailwind-variants";

export const buttonStyles = tv({
	slots: {
		base: [
			"inline-flex cursor-pointer items-center justify-center gap-2 whitespace-nowrap rounded-md font-medium text-sm transition-colors",
			"focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50",
			"[&[data-loading=true]_svg:not(:first-child)]:hidden [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0",
		],
	},
	variants: {
		variant: {
			default: "bg-gray-900 text-white shadow-sm hover:bg-gray-900/90",
			primary: "bg-primary text-white shadow-sm hover:bg-primary/90",
			secondary:
				"bg-secondary text-secondary-foreground shadow-xs hover:bg-secondary/80",
			destructive:
				"bg-destructive text-destructive-foreground shadow-xs hover:bg-destructive/90",
			outline:
				"border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground",
			ghost: "hover:bg-accent hover:text-accent-foreground",
			link: "text-primary underline-offset-4 hover:underline",
		},
		size: {
			default: "h-9 px-4 py-2",
			sm: "h-8 rounded-md px-3 text-xs",
			lg: "h-10 rounded-md px-8",
			icon: "size-9",
		},
		isLoading: {
			true: "pointer-events-none relative cursor-wait",
			false: "",
		},
	},
	compoundVariants: [
		{
			isLoading: true,
			className:
				"[&>svg]:-ml-1 [&>svg]:motion-preset-spin [&>svg]:motion-duration-1000 [&>svg]:size-4",
		},
		{
			variant: [
				"default",
				"primary",
				"secondary",
				"outline",
				"destructive",
				"ghost",
				"link",
			],
			isLoading: true,
			className: "text-white/70",
		},
		{
			size: ["icon"],
			isLoading: true,
			className: "[&>svg]:ml-0",
		},
		{
			variant: "outline",
			size: ["icon", "sm"],
			className: "border",
		},
		{
			variant: ["ghost", "link"],
			size: ["lg"],
			className: "font-semibold tracking-wide",
		},
		{
			variant: [
				"default",
				"primary",
				"secondary",
				"outline",
				"destructive",
				"ghost",
				"link",
			],
			size: ["lg"],
			className: "shadow-md",
		},
	],
	defaultVariants: {
		variant: "default",
		size: "default",
		isLoading: false,
	},
});

export type ButtonVariants = VariantProps<typeof buttonStyles>;
