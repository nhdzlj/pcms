<template>
  <div style="padding: 12px 0">
    <el-menu
      :default-active="activeMenu"
      router
      style="border-right: none"
      @select="handleMenuSelect"
    >
      <el-menu-item index="/home">
        <el-icon><HomeFilled /></el-icon>
        <span>首页</span>
      </el-menu-item>
      <el-menu-item index="/documents">
        <el-icon><Document /></el-icon>
        <span>所有文档</span>
      </el-menu-item>
      <el-menu-item index="/search">
        <el-icon><Search /></el-icon>
        <span>搜索</span>
      </el-menu-item>
    </el-menu>

    <el-divider style="margin: 8px 0" />

    <CategoryTree />
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";
import { HomeFilled, Document, Search } from "@element-plus/icons-vue";
import CategoryTree from "./CategoryTree.vue";

const route = useRoute();

const activeMenu = computed(() => {
  const path = route.path;
  if (path === "/home") return "/home";
  if (path.startsWith("/documents")) return "/documents";
  if (path.startsWith("/search")) return "/search";
  return "/home";
});

function handleMenuSelect() {
  // Handle menu selection if needed
}
</script>
