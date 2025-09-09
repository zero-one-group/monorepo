import * as Lucide from "lucide-react";
import * as ResizablePrimitive from "react-resizable-panels";
import { resizableStyles } from "./resizable.css";

interface ResizablePanelGroupProps
	extends React.ComponentProps<typeof ResizablePrimitive.PanelGroup> {
	fixed?: boolean;
}

const ResizablePanelGroup = ({
	className,
	fixed,
	...props
}: ResizablePanelGroupProps) => {
	const styles = resizableStyles();
	return (
		<ResizablePrimitive.PanelGroup
			className={
				fixed ? styles.groupFixed({ className }) : styles.group({ className })
			}
			{...props}
		/>
	);
};

const ResizablePanel = ResizablePrimitive.Panel;

const ResizableHandle = ({
	withHandle,
	className,
	...props
}: React.ComponentProps<typeof ResizablePrimitive.PanelResizeHandle> & {
	withHandle?: boolean;
}) => {
	const styles = resizableStyles();
	return (
		<ResizablePrimitive.PanelResizeHandle
			className={styles.handle({
				className: [
					styles.handleVertical(),
					styles.handleAfter(),
					styles.handleAfterVertical(),
					styles.handleRotate(),
					className,
				],
			})}
			{...props}
		>
			{withHandle && (
				<div className={styles.handleButton()}>
					<Lucide.GripVertical className={styles.handleIcon()} />
				</div>
			)}
		</ResizablePrimitive.PanelResizeHandle>
	);
};

ResizablePanelGroup.displayName = "ResizablePanelGroup";
ResizablePanel.displayName = "ResizablePanel";
ResizableHandle.displayName = "ResizableHandle";

export { ResizablePanelGroup, ResizablePanel, ResizableHandle };
