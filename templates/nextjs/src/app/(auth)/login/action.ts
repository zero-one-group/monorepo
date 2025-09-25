import consola from "consola";

export interface LoginState {
	error?: string;
	success?: boolean;
	data?: {
		email: string;
		timestamp: string;
	};
}

export const loginAction = async (
	_prevState: LoginState,
	formData: FormData,
): Promise<LoginState> => {
	await new Promise((resolve) => setTimeout(resolve, 1000));

	const email = formData.get("email");
	const password = formData.get("password");

	if (!email || !password) {
		const errorState = { error: "Email and password are required!" };
		consola.error("Login Error:", errorState);
		return errorState;
	}

	const successState = {
		success: true,
		data: {
			email: email.toString(),
			timestamp: new Date().toISOString(),
		},
	};
	consola.info("Login Success:", successState);
	return successState;
};
