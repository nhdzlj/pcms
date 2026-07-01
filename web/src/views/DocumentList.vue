<template>
  <div>
    <!-- 工具栏 -->
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px">
      <div style="display: flex; align-items: center; gap: 12px">
        <h3 style="font-size: 18px; font-weight: 600">
          {{ currentCategoryName || "所有文档" }}
        </h3>
        <el-tag v-if="currentCategoryName" closable @close="clearFilter">
          分类: {{ currentCategoryName }}
        </el-tag>
      </div>

      <div style="display: flex; gap: 8px">
        <el-select
          v-model="statusFilter"
          placeholder="状态"
          clearable
          size="default"
          style="width: 120px"
          @change="loadData"
        >
          <el-option label="全部" value="" />
          <el-option label="草稿" value="draft" />
          <el-option label="已发布" value="published" />
          <el-option label="已归档" value="archived" />
        </el-select>
        <el-button type="primary" @click="router.push('/documents/new')">
          <el-icon><Plus /></el-icon> 新建
        </el-button>
      </div>
    </div>

    <!-- 文档列表 -->
    <el-table
      :data="docStore.documents"
      v-loading="docStore.loading"
      stripe
      style="width: 100%"
      @row-click="handleRowClick"
    >
      <el-table-column prop="title" label="标题" min-width="200">
        <template #default="{ row }">
          <div style="display: flex; align-items: center; gap: 8px">
            <el-icon :size="18">
              <Document />
            </el-icon>
            <span style="font-weight: 500">{{ row.title }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="category.name" label="分类" width="120">
        <template #default="{ row }">
          <el-tag v-if="row.category" size="small" type="info" effect="plain">
            {{ row.category.name }}
          </el-tag>
          <span v-else style="color: #c0c4cc">未分类</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="90">
        <template #default="{ row }">
          <el-tag
            :type="row.status === 'published' ? 'success' : row.status === 'draft' ? 'warning' : 'info'"
            size="small"
            effect="plain"
          >
            {{ statusMap[row.status] || row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="view_count" label="浏览" width="80" align="center" />
      <el-table-column prop="updated_at" label="更新时间" width="180">
        <template #default="{ row }">
          {{ new Date(row.updated_at).toLocaleString("zh-CN") }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click.stop="router.push(`/documents/${row.id}`)">
            编辑
          </el-button>
          <el-popconfirm
            title="确定要删除这个文档吗？"
            @confirm="handleDelete(row.id)"
          >
            <template #reference>
              <el-button link type="danger" size="small" @click.stop>删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div style="margin-top: 16px; display: flex; justify-content: flex-end">
      <el-pagination
        v-model:current-page="docStore.currentPage"
        :page-size="docStore.pageSize"
        :total="docStore.total"
        layout="total, prev, pager, next"
        @current-change="loadData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Plus, Document } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { useDocumentStore } from "@/stores/document";
import { useCategoryStore } from "@/stores/category";

const route = useRoute();
const router = useRouter();
const docStore = useDocumentStore();
const catStore = useCategoryStore();

const statusFilter = ref("");

const statusMap: Record<string, string> = {
  draft: "草稿",
  published: "已发布",
  archived: "已归档",
};

const currentCategoryName = computed(() => {
  const cid = route.query.category_id;
  if (!cid) return "";
  const find = (arr: any[]): any => {
    for (const item of arr) {
      if (item.id === Number(cid)) return item;
      if (item.children) {
        const found = find(item.children);
        if (found) return found;
      }
    }
    return null;
  };
  const cat = find(catStore.tree);
  return cat?.name || "";
});

function clearFilter() {
  router.push("/documents");
}

async function loadData() {
  const categoryId = route.query.category_id
    ? Number(route.query.category_id)
    : undefined;
  await docStore.fetchList({
    page: docStore.currentPage,
    category_id: categoryId,
    status: statusFilter.value || undefined,
  });
}

function handleRowClick(row: any) {
  router.push(`/documents/${row.id}`);
}

async function handleDelete(id: number) {
  try {
    await docStore.remove(id);
    ElMessage.success("删除成功");
    loadData();
  } catch {
    ElMessage.error("删除失败");
  }
}

// URL query 变化时重新加载
watch(() => route.query.category_id, () => {
  docStore.currentPage = 1;
  loadData();
});

onMounted(() => {
  loadData();
});
</script>
