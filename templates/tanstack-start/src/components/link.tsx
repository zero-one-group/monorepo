import * as React from "react";

type LinkProps = React.ComponentPropsWithoutRef<"a"> & {
    href: string;
    newTab?: boolean;
};

const Link = React.forwardRef<HTMLAnchorElement, LinkProps>(function Link(
    props,
    ref
) {
    const { className, newTab, ...rest } = props;
    const NEW_TAB_REL = "noopener noreferrer";
    const NEW_TAB_TARGET = "_blank";
    const DEFAULT_TARGET = "_self";

    return (
        <a
            className={className}
            rel={newTab ? NEW_TAB_REL : undefined}
            target={newTab ? NEW_TAB_TARGET : DEFAULT_TARGET}
            ref={ref}
            {...rest}
        />
    );
});

export { Link };
