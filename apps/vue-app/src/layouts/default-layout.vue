<script setup lang="ts">
import AppNavbar from "#/components/app-navbar.vue";
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <AppNavbar />
    <slot />
  </div>
</template>
