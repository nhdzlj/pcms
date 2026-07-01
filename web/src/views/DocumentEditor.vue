<template>
  <div class="editor-layout" v-loading="loading">
    <!-- 顶部工具栏 -->
    <div class="editor-toolbar">
      <el-button text @click="handleBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>

      <el-input
        v-model="form.title"
        placeholder="请输入标题..."
        size="large"
        style="flex: 1; max-width: 400px"
        class="title-input"
      />

      <el-select
        v-model="form.category_id"
        placeholder="选择分类"
        clearable
        style="width: 180px"
      >
        <el-option
          v-for="cat in flatCategories"
          :key="cat.id"
          :label="'　'.repeat(cat.depth) + cat.name"
          :value="cat.id"
        />
      </el-select>

      <el-select
        v-model="form.status"
        style="width: 100px"
      >
        <el-option label="草稿" value="draft" />
        <el-option label="发布" value="published" />
        <el-option label="归档" value="archived" />
      </el-select>

      <div style="flex: 1" />

      <el-button v-if="isEdit" @click="showVersions = true">
        <el-icon><Clock /></el-icon> 版本
      </el-button>

      <el-button @click="handleSave" :loading="saving" type="primary">
        <el-icon><Check /></el-icon>
        保存
      </el-button>
    </div>

    <!-- 摘要和标签 -->
    <div class="editor-meta">
      <el-input
        v-model="form.summary"
        placeholder="文档摘要（可选）..."
        size="small"
        clearable
      />
      <el-select
        v-model="form.tag_ids"
        multiple
        filterable
        allow-create
        default-first-option
        placeholder="选择或创建标签"
        style="width: 300px"
        @create="handleCreateTag"
      >
        <el-option
          v-for="tag in allTags"
          :key="tag.id"
          :label="tag.name"
          :value="tag.id"
        />
      </el-select>
    </div>

    <!-- 编辑器 -->
    <div class="editor-body">
      <MarkdownEditor
        v-model="form.content"
        placeholder="开始写作... 支持 Markdown 格式"
        :show-split="true"
      />
    </div>

    <!-- 版本历史对话框 -->
    <el-dialog v-model="showVersions" title="版本历史" width="700px">
      <el-table :data="docStore.versions" v-loading="versionLoading" stripe max-height="400">
        <el-table-column prop="version" label="版本" width="80">
          <template #default="{ row }">
            <el-tag type="primary" effect="plain" size="small">v{{ row.version }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="150" />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString("zh-CN") }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleViewVersion(row)">
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!versionLoading && docStore.versions.length === 0" description="暂无版本记录" />
    </el-dialog>

    <!-- 版本内容对话框 -->
    <el-dialog v-model="showVersionContent" :title="`版本 v${viewingVersion?.version}`" width="800px">
      <div style="margin-bottom: 12px">
        <el-tag size="small">{{ viewingVersion?.title }}</el-tag>
        <span style="margin-left: 12px; color: var(--app-text-secondary); font-size: 13px">
          {{ viewingVersion ? new Date(viewingVersion.created_at).toLocaleString("zh-CN") : "" }}
        </span>
      </div>
      <MarkdownEditor
        :model-value="viewingVersion?.content || ''"
        :show-split="true"
        placeholder=""
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { ArrowLeft, Check, Clock } from "@element-plus/icons-vue";
import { useDocumentStore } from "@/stores/document";
import { useCategoryStore } from "@/stores/category";
import { getTags, createTag } from "@/api/tag";
import type { Tag } from "@/api/tag";
import type { DocumentVersion } from "@/api/document";
import MarkdownEditor from "@/components/MarkdownEditor.vue";

const route = useRoute();
const router = useRouter();
const docStore = useDocumentStore();
const catStore = useCategoryStore();

const loading = ref(false);
const saving = ref(false);
const isEdit = computed(() => !!route.params.id);

const allTags = ref<Tag[]>([]);
const showVersions = ref(false);
const showVersionContent = ref(false);
const versionLoading = ref(false);
const viewingVersion = ref<DocumentVersion | null>(null);

interface FlatCategory {
  id: number;
  name: string;
  depth: number;
}

function flattenTree(
  categories: any[],
  depth = 0
): FlatCategory[] {
  const result: FlatCategory[] = [];
  for (const cat of categories) {
    result.push({ id: cat.id, name: cat.name, depth });
    if (cat.children) {
      result.push(...flattenTree(cat.children, depth + 1));
    }
  }
  return result;
}

const flatCategories = computed(() => flattenTree(catStore.tree));

const form = reactive({
  title: "",
  content: "",
  category_id: null as number | null,
  status: "draft",
  summary: "",
  tag_ids: [] as number[],
});

onMounted(async () => {
  await catStore.fetchTree();
  await loadTags();
  if (isEdit.value) {
    const id = Number(route.params.id);
    loading.value = true;
    try {
      await docStore.fetchById(id);
      if (docStore.current) {
        form.title = docStore.current.title;
        form.content = docStore.current.content;
        form.category_id = docStore.current.category_id;
        form.status = docStore.current.status;
        form.summary = docStore.current.summary || "";
        form.tag_ids = docStore.current.tags?.map((t: any) => t.id) || [];
      }
    } catch {
      ElMessage.error("文档不存在");
      router.push("/documents");
    } finally {
      loading.value = false;
    }
  }
});

async function loadTags() {
  try {
    allTags.value = await getTags();
  } catch {
    // ignore
  }
}

async function handleCreateTag(name: string) {
  try {
    const tag = await createTag(name);
    allTags.value.push(tag);
    // 选中新创建的标签
    if (tag.id) {
      form.tag_ids = [...form.tag_ids, tag.id];
      await nextTick();
      // Force refresh the options list
      allTags.value = [...allTags.value];
    }
    ElMessage.success("标签已创建");
  } catch {
    ElMessage.error("创建标签失败");
  }
}

async function handleSave() {
  if (!form.title.trim()) {
    ElMessage.warning("请输入标题");
    return;
  }

  saving.value = true;
  try {
    if (isEdit.value) {
      const id = Number(route.params.id);
      await docStore.update(id, {
        title: form.title,
        content: form.content,
        category_id: form.category_id,
        status: form.status,
        summary: form.summary,
        tag_ids: form.tag_ids,
      });
      ElMessage.success("保存成功");
    } else {
      const doc = await docStore.create({
        title: form.title,
        content: form.content,
        category_id: form.category_id,
        status: form.status,
        summary: form.summary,
        tag_ids: form.tag_ids,
      });
      ElMessage.success("创建成功");
      router.replace(`/documents/${doc.id}`);
    }
  } catch (err: any) {
    ElMessage.error(err.message || "保存失败");
  } finally {
    saving.value = false;
  }
}

function handleBack() {
  router.push("/documents");
}

watch(showVersions, async (val) => {
  if (val && isEdit.value) {
    versionLoading.value = true;
    try {
      await docStore.fetchVersions(Number(route.params.id));
    } finally {
      versionLoading.value = false;
    }
  }
});

async function handleViewVersion(version: DocumentVersion) {
  viewingVersion.value = version;
  showVersionContent.value = true;
}
</script>

<style scoped>
.title-input :deep(.el-input__inner) {
  font-size: 18px;
  font-weight: 600;
  border: none;
  padding-left: 0;
}
.title-input :deep(.el-input__wrapper) {
  box-shadow: none !important;
}

.editor-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  border-bottom: 1px solid var(--app-border);
  background: var(--app-header-bg);
  flex-shrink: 0;
}
.editor-meta .el-input {
  flex: 1;
  max-width: 400px;
}
</style>
