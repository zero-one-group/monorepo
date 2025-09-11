/**
 * A custom Link component that extends the functionality of the React Router Link component.
 * It adds support for opening links in a new tab and applies a consistent set of styles.
 *
 * @param props - The props for the Link component, including the standard Link props from React Router.
 * @param props.newTab - Whether to open the link in a new tab.
 * @param props.className - Additional CSS classes to apply to the link.
 * @param ref - A ref to the underlying anchor element.
 * @returns A React Router Link component with the custom functionality and styles applied.
 *
 * Example usage:
 * ```tsx
 * <Link href="/path" className="custom-class">Link Text</Link>
 * ```
 */

import * as React from "react";
import {
	Link as RouterLink,
	type LinkProps as RouterLinkProps,
} from "react-router";

interface LinkProps extends Omit<RouterLinkProps, "to"> {
	href: string;
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
			to={props.href}
			className={className}
			rel={newTab ? NEW_TAB_REL : undefined}
			target={newTab ? NEW_TAB_TARGET : DEFAULT_TARGET}
			ref={ref}
			{...rest}
		/>
	);
});

export { Link };
