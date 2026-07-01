<template>
  <div>
    <!-- 搜索栏 -->
    <div style="margin-bottom: 24px">
      <el-input
        v-model="keyword"
        placeholder="搜索文档标题和内容..."
        size="large"
        clearable
        @keyup.enter="handleSearch"
        @clear="handleClear"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
        <template #append>
          <el-button
            type="primary"
            @click="handleSearch"
            :loading="loading"
          >
            搜索
          </el-button>
        </template>
      </el-input>
    </div>

    <!-- 搜索结果提示 -->
    <div v-if="searched && keyword" style="margin-bottom: 16px; color: var(--app-text-secondary)">
      找到 <strong>{{ docStore.total }}</strong> 条与 "<strong>{{ searchedKeyword }}</strong>" 相关的结果
    </div>

    <!-- 结果列表 -->
    <el-table
      v-if="searched"
      :data="docStore.documents"
      v-loading="loading"
      stripe
      style="width: 100%"
      @row-click="handleRowClick"
    >
      <el-table-column label="标题" min-width="250">
        <template #default="{ row }">
          <div style="display: flex; align-items: center; gap: 8px">
            <el-icon :size="18"><Document /></el-icon>
            <span style="font-weight: 500; cursor: pointer; color: var(--app-primary)">
              {{ row.title }}
            </span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="summary" label="摘要" min-width="300">
        <template #default="{ row }">
          <span style="color: var(--app-text-secondary); font-size: 13px">
            {{ row.summary || truncate(row.content, 100) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="category.name" label="分类" width="120">
        <template #default="{ row }">
          <el-tag v-if="row.category" size="small" type="info" effect="plain">
            {{ row.category.name }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="更新时间" width="180">
        <template #default="{ row }">
          {{ new Date(row.updated_at).toLocaleString("zh-CN") }}
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div
      v-if="searched && docStore.total > 0"
      style="margin-top: 16px; display: flex; justify-content: flex-end"
    >
      <el-pagination
        v-model:current-page="docStore.currentPage"
        :page-size="docStore.pageSize"
        :total="docStore.total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>

    <!-- 空状态 -->
    <el-empty v-if="searched && docStore.documents.length === 0 && !loading" description="未找到相关文档" />

    <div v-if="!searched" style="text-align: center; padding: 80px 0; color: var(--app-text-secondary)">
      <el-icon :size="64" style="color: #c0c4cc; margin-bottom: 16px">
        <Search />
      </el-icon>
      <div style="font-size: 16px">输入关键词搜索您的知识库</div>
      <div style="font-size: 13px; margin-top: 8px">支持搜索标题和内容</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { Search, Document } from "@element-plus/icons-vue";
import { useDocumentStore } from "@/stores/document";

const router = useRouter();
const docStore = useDocumentStore();

const keyword = ref("");
const searched = ref(false);
const searchedKeyword = ref("");
const loading = ref(false);

async function handleSearch() {
  const kw = keyword.value.trim();
  if (!kw) {
    ElMessage.warning("请输入搜索关键词");
    return;
  }

  loading.value = true;
  searched.value = true;
  searchedKeyword.value = kw;
  try {
    await docStore.search(kw, 1);
  } finally {
    loading.value = false;
  }
}

function handleClear() {
  searched.value = false;
  searchedKeyword.value = "";
  keyword.value = "";
}

function handlePageChange(page: number) {
  docStore.search(searchedKeyword.value, page);
}

function handleRowClick(row: any) {
  router.push(`/documents/${row.id}`);
}

function truncate(text: string, len: number): string {
  if (!text) return "";
  // Strip markdown syntax roughly
  let plain = text
    .replace(/[#*`~\[\]()>]/g, "")
    .replace(/\n/g, " ")
    .trim();
  return plain.length > len ? plain.slice(0, len) + "..." : plain;
}
</script>
