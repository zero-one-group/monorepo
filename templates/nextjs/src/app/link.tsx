/**
 * A custom Link component that extends the functionality of the Next.js Link component.
 * It adds support for opening links in a new tab and applies a consistent set of styles.
 *
 * @param props - The props for the Link component, including the standard Link props from Next.js.
 * @param props.newTab - Whether to open the link in a new tab.
 * @param props.className - Additional CSS classes to apply to the link.
 * @param ref - A ref to the underlying anchor element.
 * @returns A Next.js Link component with the custom functionality and styles applied.
 *
 * Example usage:
 * ```tsx
 * <Link href="/path" className="custom-class">Link Text</Link>
 * ```
 */

import type { LinkProps as RouterLinkProps } from "next/link";
import { default as RouterLink } from "next/link";
import * as React from "react";
import { clx } from "#/libs/utils";

interface LinkProps extends Omit<RouterLinkProps, "to"> {
	newTab?: boolean;
}

const Link = React.forwardRef(function Component(
	props: LinkProps & React.ComponentPropsWithoutRef<"a">,
	ref: React.ForwardedRef<HTMLAnchorElement>,
) {
	const { className, newTab, ...rest } = props;
	const NEW_TAB_REL = "noopener noreferrer";
	const NEW_TAB_TARGET = "_blank";
	const DEFAULT_TARGET = "_self";

	return (
		<RouterLink
			className={clx("text-inherit dark:text-inherit", className)}
			target={newTab ? NEW_TAB_TARGET : DEFAULT_TARGET}
			rel={newTab ? NEW_TAB_REL : undefined}
			ref={ref}
			{...rest}
		/>
	);
});

export default Link;
