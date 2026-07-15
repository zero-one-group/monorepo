<script setup lang="ts">
import { useRouter } from "vue-router";

const router = useRouter();
</script>

<template>
  <main class="flex flex-1 flex-col items-center justify-center gap-4">
    <h1 class="text-4xl font-bold">404</h1>
    <p class="text-gray-500">Page not found</p>
    <button class="text-blue-600 underline" @click="router.push('/')">
      Go home
    </button>
  </main>
</template>
