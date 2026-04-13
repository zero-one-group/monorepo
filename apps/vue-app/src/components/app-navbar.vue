<script setup lang="ts">
import { RouterLink } from "vue-router";

const appName = "vue-app";
</script>

<template>
  <nav class="flex items-center justify-between border-b px-6 py-3">
    <RouterLink class="text-lg font-semibold" to="/">
      {{ appName }}
    </RouterLink>
    <div class="flex items-center gap-6">
      <RouterLink
        class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-100"
        exact-active-class="font-medium text-gray-900 dark:text-gray-100"
        to="/"
      >
        Home
      </RouterLink>
      <!-- active-class (not exact) so sub-routes like /login?redirect=... also match -->
      <RouterLink
        class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-100"
        active-class="font-medium text-gray-900 dark:text-gray-100"
        to="/login"
      >
        Login
      </RouterLink>
    </div>
  </nav>
</template>
