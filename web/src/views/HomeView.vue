<template>
  <div>
    <el-row :gutter="16" style="margin-bottom: 16px">
      <el-col :span="6" v-for="stat in stats" :key="stat.label">
        <el-card shadow="hover">
          <div style="text-align: center">
            <div style="font-size: 28px; font-weight: 700; color: var(--app-primary)">
              {{ stat.value }}
            </div>
            <div style="font-size: 13px; color: var(--app-text-secondary); margin-top: 4px">
              {{ stat.label }}
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px">
      <h3 style="font-size: 18px; font-weight: 600">最近文档</h3>
      <el-button type="primary" @click="router.push('/documents/new')">
        <el-icon><Plus /></el-icon> 新建文档
      </el-button>
    </div>

    <el-row :gutter="16" v-loading="loading">
      <el-col
        :span="8"
        v-for="doc in recentDocs"
        :key="doc.id"
        style="margin-bottom: 16px"
      >
        <DocumentCard :doc="doc" />
      </el-col>
    </el-row>

    <el-empty v-if="!loading && recentDocs.length === 0" description="还没有文档，点击上方按钮创建" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { Plus } from "@element-plus/icons-vue";
import { useDocumentStore } from "@/stores/document";
import { useCategoryStore } from "@/stores/category";
import DocumentCard from "@/components/DocumentCard.vue";

const router = useRouter();
const docStore = useDocumentStore();
const catStore = useCategoryStore();
const loading = ref(false);

function countCategories(tree: any[]): number {
  let count = 0;
  for (const item of tree) {
    count++;
    if (item.children) {
      count += countCategories(item.children);
    }
  }
  return count;
}

const stats = computed(() => [
  { label: "总文档", value: docStore.total },
  { label: "最近更新", value: docStore.documents.length },
  { label: "分类数", value: countCategories(catStore.tree) },
]);

const recentDocs = computed(() => docStore.documents.slice(0, 12));

onMounted(async () => {
  loading.value = true;
  try {
    await Promise.all([docStore.fetchList({ page: 1 }), catStore.fetchTree()]);
  } finally {
    loading.value = false;
  }
});
</script>
