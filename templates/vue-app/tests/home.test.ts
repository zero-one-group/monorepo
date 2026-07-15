import { mount } from "@vue/test-utils";
import { describe, expect, it } from "vitest";
import { createMemoryHistory, createRouter } from "vue-router";
import App from "../src/App.vue";
import HomeView from "../src/views/home.vue";

function makeRouter() {
	return createRouter({
		history: createMemoryHistory(),
		routes: [{ path: "/", component: HomeView }],
	});
}

describe("App", () => {
	it("renders a nav element", async () => {
		const router = makeRouter();
		const wrapper = mount(App, { global: { plugins: [router] } });
		await router.isReady();
		expect(wrapper.find("nav").exists()).toBe(true);
	});

	it("renders Home and Login nav links", async () => {
		const router = makeRouter();
		const wrapper = mount(App, { global: { plugins: [router] } });
		await router.isReady();
		const hrefs = wrapper.findAll("a").map((l) => l.attributes("href"));
		expect(hrefs).toContain("/");
		expect(hrefs).toContain("/login");
	});

	it("renders an h1 on the home route", async () => {
		const router = makeRouter();
		await router.push("/");
		const wrapper = mount(App, { global: { plugins: [router] } });
		await router.isReady();
		expect(wrapper.find("h1").exists()).toBe(true);
	});
});
