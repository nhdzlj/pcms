<template>
  <div class="app-layout">
    <AppHeader />
    <div class="app-body">
      <aside class="app-sidebar">
        <AppSidebar />
      </aside>
      <main class="app-main">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppHeader from "@/components/AppHeader.vue";
import AppSidebar from "@/components/AppSidebar.vue";
</script>
