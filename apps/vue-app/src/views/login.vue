<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const email = ref("");
const password = ref("");

async function handleSubmit() {
	// TODO: replace with real auth
	await router.push("/");
}
</script>

<template>
  <main class="flex flex-1 flex-col items-center justify-center">
    <div class="w-full max-w-sm space-y-4 rounded-lg border p-6">
      <h1 class="text-2xl font-bold">Login</h1>
      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="space-y-1">
          <label class="text-sm font-medium" for="email">Email</label>
          <input
            id="email"
            v-model="email"
            class="w-full rounded-md border px-3 py-2 text-sm"
            placeholder="you@example.com"
            required
            type="email"
          />
        </div>
        <div class="space-y-1">
          <label class="text-sm font-medium" for="password">Password</label>
          <input
            id="password"
            v-model="password"
            class="w-full rounded-md border px-3 py-2 text-sm"
            placeholder="••••••••"
            required
            type="password"
          />
        </div>
        <button
          class="w-full rounded-md bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700 dark:bg-gray-100 dark:text-gray-900 dark:hover:bg-gray-200"
          type="submit"
        >
          Sign in
        </button>
      </form>
    </div>
  </main>
</template>
