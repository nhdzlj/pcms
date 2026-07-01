<template>
  <div class="doc-card el-card" @click="handleClick">
    <div class="doc-card-title">{{ doc.title }}</div>
    <div v-if="doc.summary" class="doc-card-summary">{{ doc.summary }}</div>
    <div v-else class="doc-card-summary" style="color: #c0c4cc">
      暂无摘要
    </div>
    <div class="doc-card-meta">
      <el-tag
        v-if="doc.category"
        size="small"
        type="info"
        effect="plain"
      >
        {{ doc.category.name }}
      </el-tag>
      <el-tag
        v-if="doc.status === 'draft'"
        size="small"
        type="warning"
        effect="plain"
      >
        草稿
      </el-tag>
      <el-tag
        v-else-if="doc.status === 'published'"
        size="small"
        type="success"
        effect="plain"
      >
        已发布
      </el-tag>
      <span style="margin-left: auto">
        {{ formatDate(doc.updated_at) }}
      </span>
      <span>
        <el-icon><View /></el-icon>
        {{ doc.view_count }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { View } from "@element-plus/icons-vue";
import type { Document } from "@/api/document";

const props = defineProps<{
  doc: Document;
}>();

const router = useRouter();

function handleClick() {
  router.push(`/documents/${props.doc.id}`);
}

function formatDate(dateStr: string): string {
  if (!dateStr) return "";
  const d = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - d.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (days === 0) return "今天";
  if (days === 1) return "昨天";
  if (days < 7) return `${days}天前`;
  return d.toLocaleDateString("zh-CN");
}
</script>
