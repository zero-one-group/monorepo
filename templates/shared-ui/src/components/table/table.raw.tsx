import * as React from "react";
import { tableStyles } from "./table.css";

const Table = React.forwardRef<
	HTMLTableElement,
	React.HTMLAttributes<HTMLTableElement>
>(({ className, ...props }, ref) => {
	const styles = tableStyles();
	return (
		<div className={styles.wrapper()}>
			<table ref={ref} className={styles.table({ className })} {...props} />
		</div>
	);
});

const TableHeader = React.forwardRef<
	HTMLTableSectionElement,
	React.HTMLAttributes<HTMLTableSectionElement>
>(({ className, ...props }, ref) => {
	const styles = tableStyles();
	return (
		<thead ref={ref} className={styles.header({ className })} {...props} />
	);
});

const TableBody = React.forwardRef<
	HTMLTableSectionElement,
	React.HTMLAttributes<HTMLTableSectionElement>
>(({ className, ...props }, ref) => {
	const styles = tableStyles();
	return <tbody ref={ref} className={styles.body({ className })} {...props} />;
});

const TableFooter = React.forwardRef<
	HTMLTableSectionElement,
	React.HTMLAttributes<HTMLTableSectionElement>
>(({ className, ...props }, ref) => {
	const styles = tableStyles();
	return (
		<tfoot ref={ref} className={styles.footer({ className })} {...props} />
	);
});

const TableRow = React.forwardRef<
	HTMLTableRowElement,
	React.HTMLAttributes<HTMLTableRowElement>
>(({ className, ...props }, ref) => {
	const styles = tableStyles();
	return <tr ref={ref} className={styles.row({ className })} {...props} />;
});

const TableHead = React.forwardRef<
	HTMLTableCellElement,
	React.ThHTMLAttributes<HTMLTableCellElement>
>(({ className, ...props }, ref) => {
	const styles = tableStyles();
	return <th ref={ref} className={styles.head({ className })} {...props} />;
});

const TableCell = React.forwardRef<
	HTMLTableCellElement,
	React.TdHTMLAttributes<HTMLTableCellElement>
>(({ className, ...props }, ref) => {
	const styles = tableStyles();
	return <td ref={ref} className={styles.cell({ className })} {...props} />;
});

const TableCaption = React.forwardRef<
	HTMLTableCaptionElement,
	React.HTMLAttributes<HTMLTableCaptionElement>
>(({ className, ...props }, ref) => {
	const styles = tableStyles();
	return (
		<caption ref={ref} className={styles.caption({ className })} {...props} />
	);
});

Table.displayName = "Table";
TableHeader.displayName = "TableHeader";
TableBody.displayName = "TableBody";
TableFooter.displayName = "TableFooter";
TableRow.displayName = "TableRow";
TableHead.displayName = "TableHead";
TableCell.displayName = "TableCell";
TableCaption.displayName = "TableCaption";

export {
	Table,
	TableHeader,
	TableBody,
	TableFooter,
	TableHead,
	TableRow,
	TableCell,
	TableCaption,
};
