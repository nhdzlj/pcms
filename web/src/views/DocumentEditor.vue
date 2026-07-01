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

      <el-button @click="handleSave" :loading="saving" type="primary">
        <el-icon><Check /></el-icon>
        保存
      </el-button>
    </div>

    <!-- 编辑器 -->
    <div class="editor-body">
      <MarkdownEditor
        v-model="form.content"
        placeholder="开始写作... 支持 Markdown 格式"
        :show-split="true"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { ArrowLeft, Check } from "@element-plus/icons-vue";
import { useDocumentStore } from "@/stores/document";
import { useCategoryStore } from "@/stores/category";
import MarkdownEditor from "@/components/MarkdownEditor.vue";

const route = useRoute();
const router = useRouter();
const docStore = useDocumentStore();
const catStore = useCategoryStore();

const loading = ref(false);
const saving = ref(false);
const isEdit = computed(() => !!route.params.id);

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
</style>
