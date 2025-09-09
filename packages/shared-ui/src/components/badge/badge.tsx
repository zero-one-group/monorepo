import type * as React from "react";
import { type BadgeVariants, badgeStyles } from "./badge.css";

export interface BadgeProps
	extends React.HTMLAttributes<HTMLDivElement>,
		BadgeVariants {}

function Badge({ className, variant, size, rounded, ...props }: BadgeProps) {
	const styles = badgeStyles({ variant, size, rounded });
	return <div className={styles.base({ className })} {...props} />;
}

Badge.displayName = "Badge";

export { Badge };
