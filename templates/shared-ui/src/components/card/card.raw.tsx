import * as React from "react";
import { cardStyles } from "./card.css";

const Card = React.forwardRef<
	HTMLDivElement,
	React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => {
	const styles = cardStyles();
	return <div ref={ref} className={styles.base({ className })} {...props} />;
});

const CardHeader = React.forwardRef<
	HTMLDivElement,
	React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => {
	const styles = cardStyles();
	return <div ref={ref} className={styles.header({ className })} {...props} />;
});

const CardTitle = React.forwardRef<
	HTMLParagraphElement,
	React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => {
	const styles = cardStyles();
	return <h3 ref={ref} className={styles.title({ className })} {...props} />;
});

const CardDescription = React.forwardRef<
	HTMLParagraphElement,
	React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => {
	const styles = cardStyles();
	return (
		<p ref={ref} className={styles.description({ className })} {...props} />
	);
});

const CardContent = React.forwardRef<
	HTMLDivElement,
	React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => {
	const styles = cardStyles();
	return <div ref={ref} className={styles.content({ className })} {...props} />;
});

const CardFooter = React.forwardRef<
	HTMLDivElement,
	React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => {
	const styles = cardStyles();
	return <div ref={ref} className={styles.footer({ className })} {...props} />;
});

Card.displayName = "Card";
CardHeader.displayName = "CardHeader";
CardTitle.displayName = "CardTitle";
CardDescription.displayName = "CardDescription";
CardContent.displayName = "CardContent";
CardFooter.displayName = "CardFooter";

export {
	Card,
	CardHeader,
	CardFooter,
	CardTitle,
	CardDescription,
	CardContent,
};
