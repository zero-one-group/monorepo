import { tv } from "tailwind-variants";

export const errorStyles = tv({
	slots: {
		wrapper: "relative min-h-screen overflow-hidden bg-background",
		decorativeGradient: "absolute inset-0 overflow-hidden",
		gradientInner: "-inset-[10px] absolute opacity-50",
		gradientBg: [
			"absolute top-0 h-[40rem] w-full",
			"bg-gradient-to-b from-muted via-transparent to-transparent",
			"before:absolute before:inset-0 before:bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] before:from-muted/10 before:via-transparent before:to-transparent",
			"after:absolute after:inset-0 after:bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] after:from-primary/10 after:via-transparent after:to-transparent",
		],
		content:
			"relative flex min-h-screen flex-col items-center justify-center px-4 py-16 sm:px-6 lg:px-8",
		container: "relative z-20 text-center",
		errorCode: "font-bold text-2xl text-error",
		title: "mt-4 font-bold text-3xl text-foreground tracking-tight sm:text-5xl",
		description: "mt-6 text-base text-muted-foreground leading-7",
		pre: "w-full overflow-x-auto p-4",
		code: "text-muted-foreground",
		actions: "mt-10 flex items-center justify-center gap-x-4",
		primaryButton: [
			"min-w-[140px] cursor-pointer rounded-md bg-primary px-4 py-2.5 font-semibold text-sm text-primary-foreground",
			"transition-all duration-200 hover:bg-primary/90 hover:shadow-primary/20 hover:shadow-lg",
			"focus:outline-none focus:ring-2 focus:ring-ring/50 focus:ring-offset-2",
			"focus:ring-offset-background",
			"border border-border",
		],
		secondaryButton: [
			"min-w-[140px] rounded-md border border-border bg-background/80 px-4 py-2.5",
			"cursor-pointer font-semibold text-foreground text-sm",
			"transition-all duration-200 hover:border-border/30 hover:bg-accent",
			"hover:text-accent-foreground hover:shadow-ring/10 hover:shadow-lg",
			"focus:outline-none focus:ring-2 focus:ring-ring/50 focus:ring-offset-2",
			"focus:ring-offset-background",
		],
		decorativeCode:
			"pointer-events-none fixed inset-0 z-10 flex select-none items-center justify-center",
		decorativeText:
			"font-black text-[12rem] text-error/10 mix-blend-overlay sm:text-[16rem] md:text-[20rem]",
	},
});
