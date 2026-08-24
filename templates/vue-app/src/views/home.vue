<script setup lang="ts">
const title = "{{ package_name | kebab_case }}";
</script>

<template>
  <main class="flex flex-1 flex-col items-center justify-center gap-4">
    <h1 class="text-4xl font-bold">{{ "{{" }} title {{ "}}" }}</h1>
    <p class="text-sm text-gray-500">
      Get started by editing
      <code class="rounded bg-gray-100 px-1 py-0.5 font-mono text-xs dark:bg-gray-800">
        src/views/home.vue
      </code>
    </p>
  </main>
</template>
