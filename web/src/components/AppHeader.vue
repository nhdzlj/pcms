<template>
  <header class="app-header">
    <div style="display: flex; align-items: center; gap: 12px; flex: 1">
      <el-icon :size="24" style="color: var(--app-primary)"><Monitor /></el-icon>
      <span style="font-size: 18px; font-weight: 700">PCMS</span>
      <el-divider direction="vertical" />
      <span style="font-size: 13px; color: var(--app-text-secondary)">
        个人认知管理系统
      </span>
    </div>

    <div style="display: flex; align-items: center; gap: 12px">
      <el-button @click="router.push('/search')">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>

      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新建文档
      </el-button>

      <el-dropdown trigger="click">
        <div style="display: flex; align-items: center; gap: 6px; cursor: pointer">
          <el-avatar :size="32">
            {{ authStore.username?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <span style="font-size: 14px">{{ authStore.username }}</span>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="authStore.logout()">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { Monitor, Search, Plus, SwitchButton } from "@element-plus/icons-vue";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const authStore = useAuthStore();

function handleCreate() {
  router.push("/documents/new");
}
</script>
