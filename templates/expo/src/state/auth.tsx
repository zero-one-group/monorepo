import React, {
	createContext,
	useCallback,
	useContext,
	useMemo,
	useReducer,
} from "react";

export type User = {
	id: string;
	email: string;
};

type AuthStatus = "signedOut" | "signingIn" | "signedIn";

type State = {
	user: User | null;
	status: AuthStatus;
	error: string | null;
};

type Action =
	| { type: "SIGN_IN_REQUEST" }
	| { type: "SIGN_IN_SUCCESS"; user: User }
	| { type: "SIGN_IN_FAILURE"; error: string }
	| { type: "SIGN_OUT" }
	| { type: "CLEAR_ERROR" };

function reducer(state: State, action: Action): State {
	switch (action.type) {
		case "SIGN_IN_REQUEST":
			return { ...state, status: "signingIn", error: null };
		case "SIGN_IN_SUCCESS":
			return { user: action.user, status: "signedIn", error: null };
		case "SIGN_IN_FAILURE":
			return { ...state, status: "signedOut", error: action.error };
		case "SIGN_OUT":
			return { user: null, status: "signedOut", error: null };
		case "CLEAR_ERROR":
			return { ...state, error: null };
		default:
			return state;
	}
}

export type AuthContextValue = {
	user: User | null;
	status: AuthStatus;
	error: string | null;
	signIn: (params: { email: string; password: string }) => Promise<void>;
	signOut: () => void;
	clearError: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

const initialState: State = {
	user: null,
	status: "signedOut",
	error: null,
};

export function AuthProvider({
	children,
}: React.PropsWithChildren): React.ReactElement {
	const [state, dispatch] = useReducer(reducer, initialState);

	const clearError = useCallback(() => {
		dispatch({ type: "CLEAR_ERROR" });
	}, []);

	const signOut = useCallback(() => {
		dispatch({ type: "SIGN_OUT" });
	}, []);

	const signIn = useCallback(
		async ({ email, password }: { email: string; password: string }) => {
			dispatch({ type: "SIGN_IN_REQUEST" });

			const normalizedEmail = email.trim().toLowerCase();
			if (!normalizedEmail) {
				const error = "Email is required.";
				dispatch({ type: "SIGN_IN_FAILURE", error });
				throw new Error(error);
			}
			if (!normalizedEmail.includes("@")) {
				const error = "Enter a valid email.";
				dispatch({ type: "SIGN_IN_FAILURE", error });
				throw new Error(error);
			}
			if (password.length < 4) {
				const error = "Password must be at least 4 characters.";
				dispatch({ type: "SIGN_IN_FAILURE", error });
				throw new Error(error);
			}

			// Simulate an async call (API / secure storage / etc.)
			await Promise.resolve();

			dispatch({
				type: "SIGN_IN_SUCCESS",
				user: { id: "user_1", email: normalizedEmail },
			});
		},
		[],
	);

	const value = useMemo<AuthContextValue>(
		() => ({
			user: state.user,
			status: state.status,
			error: state.error,
			signIn,
			signOut,
			clearError,
		}),
		[state.user, state.status, state.error, signIn, signOut, clearError],
	);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
	const value = useContext(AuthContext);
	if (!value) {
		throw new Error("useAuth must be used within an AuthProvider");
	}
	return value;
}
