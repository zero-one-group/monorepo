import { flushPromises, mount } from "@vue/test-utils";
import { describe, expect, it } from "vitest";
import { createMemoryHistory, createRouter } from "vue-router";
import App from "../src/App.vue";
import HomeView from "../src/views/home.vue";
import LoginView from "../src/views/login.vue";
import NotFoundView from "../src/views/not-found.vue";

function makeRouter() {
	return createRouter({
		history: createMemoryHistory(),
		routes: [
			{ path: "/", component: HomeView },
			{ path: "/login", component: LoginView },
			{ path: "/:pathMatch(.*)*", component: NotFoundView },
		],
	});
}

describe("LoginView", () => {
	it("renders email and password inputs", async () => {
		const router = makeRouter();
		await router.push("/login");
		const wrapper = mount(App, { global: { plugins: [router] } });
		await router.isReady();
		expect(wrapper.find('input[type="email"]').exists()).toBe(true);
		expect(wrapper.find('input[type="password"]').exists()).toBe(true);
	});

	it("renders a submit button", async () => {
		const router = makeRouter();
		await router.push("/login");
		const wrapper = mount(App, { global: { plugins: [router] } });
		await router.isReady();
		expect(wrapper.find('button[type="submit"]').exists()).toBe(true);
	});

	it("navigates to home on form submit", async () => {
		const router = makeRouter();
		await router.push("/login");
		const wrapper = mount(App, { global: { plugins: [router] } });
		await router.isReady();
		await wrapper.find("form").trigger("submit");
		await flushPromises();
		expect(router.currentRoute.value.path).toBe("/");
	});
});
