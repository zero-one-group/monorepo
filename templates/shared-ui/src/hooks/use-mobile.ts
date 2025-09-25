import { useEffect, useState } from "react";

/**
 * React hook to detect if viewport is in mobile size
 * @param breakpoint - Mobile breakpoint in pixels (default: 768)
 * @returns boolean indicating if viewport is in mobile size
 */
export const useIsMobile = (breakpoint = 768): boolean => {
	const [isMobile, setIsMobile] = useState<boolean>(false);

	useEffect(() => {
		const checkMobile = () => {
			setIsMobile(window.innerWidth < breakpoint);
		};

		// Initial check
		checkMobile();

		// Add resize listener
		window.addEventListener("resize", checkMobile);

		// Cleanup
		return () => {
			window.removeEventListener("resize", checkMobile);
		};
	}, [breakpoint]);

	return isMobile;
};
