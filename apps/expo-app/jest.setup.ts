import "@testing-library/react-native/matchers";

jest.mock("react-native-reanimated", () =>
	require("react-native-reanimated/mock"),
);
