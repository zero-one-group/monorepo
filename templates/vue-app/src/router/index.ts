import { createRouter, createWebHistory } from "vue-router";
import HomeView from "#/views/home.vue";

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: "/", component: HomeView },
		{ path: "/login", component: () => import("#/views/login.vue") },
		{
			path: "/:pathMatch(.*)*",
			component: () => import("#/views/not-found.vue"),
		},
	],
});

export default router;
