<template>
  <div class="category-tree">
    <div class="tree-header">
      <span>分类目录</span>
      <el-button :icon="Plus" size="small" circle @click="handleAdd()" />
    </div>

    <el-tree
      ref="treeRef"
      :data="store.tree"
      :props="{ children: 'children', label: 'name' }"
      node-key="id"
      :expand-on-click-node="false"
      :filter-node-method="() => true"
      highlight-current
      :current-node-key="store.selectedId"
      @node-click="handleNodeClick"
      @node-contextmenu="handleContextMenu"
      draggable
      :allow-drop="allowDrop"
      :allow-drag="() => true"
      @node-drop="handleDrop"
    >
      <template #default="{ node, data }">
        <div class="tree-node-content" @contextmenu.prevent="handleContextMenu($event, data)">
          <el-icon style="margin-right: 6px; font-size: 16px">
            <Folder v-if="node.childNodes.length > 0 || !node.parent" />
            <Document v-else />
          </el-icon>
          <span class="tree-node-label">{{ node.label }}</span>
          <el-dropdown
            trigger="click"
            @command="(cmd: string) => handleCommand(cmd, data)"
          >
            <el-button
              :size="'small'"
              :icon="MoreFilled"
              text
              style="margin-left: auto; opacity: 0.6"
              @click.stop
            />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="add">
                  <el-icon><Plus /></el-icon>添加子分类
                </el-dropdown-item>
                <el-dropdown-item command="edit">
                  <el-icon><Edit /></el-icon>重命名
                </el-dropdown-item>
                <el-dropdown-item command="delete" divided>
                  <el-icon><Delete /></el-icon>
                  <span style="color: var(--app-danger)">删除</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </template>
    </el-tree>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingCategory ? '编辑分类' : '新建分类'"
      width="400px"
    >
      <el-form @submit.prevent="handleSave">
        <el-form-item label="名称">
          <el-input
            v-model="formName"
            placeholder="分类名称"
            maxlength="128"
            ref="nameInput"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  Folder,
  Document,
  Plus,
  MoreFilled,
  Edit,
  Delete,
} from "@element-plus/icons-vue";
import { useCategoryStore } from "@/stores/category";
import type { Category } from "@/api/category";

const store = useCategoryStore();
const router = useRouter();

const treeRef = ref();

// 对话框
const dialogVisible = ref(false);
const formName = ref("");
const saving = ref(false);
const editingCategory = ref<Category | null>(null);
const parentCategory = ref<Category | null>(null);
const nameInput = ref();

function handleNodeClick(data: Category) {
  store.select(data.id);
  router.push({
    path: "/documents",
    query: { category_id: data.id },
  });
}

function handleContextMenu(event: Event, data: Category) {
  // Context menu triggers handled by dropdown
}

function handleAdd(parent?: Category) {
  editingCategory.value = null;
  parentCategory.value = parent || null;
  formName.value = "";
  dialogVisible.value = true;
  nextTick(() => nameInput.value?.focus());
}

function handleCommand(cmd: string, data: Category) {
  switch (cmd) {
    case "add":
      handleAdd(data);
      break;
    case "edit":
      editingCategory.value = data;
      formName.value = data.name;
      dialogVisible.value = true;
      nextTick(() => nameInput.value?.focus());
      break;
    case "delete":
      handleDelete(data);
      break;
  }
}

async function handleSave() {
  if (!formName.value.trim()) {
    ElMessage.warning("请输入分类名称");
    return;
  }

  saving.value = true;
  try {
    if (editingCategory.value) {
      await store.update(editingCategory.value.id, {
        name: formName.value.trim(),
      });
      ElMessage.success("更新成功");
    } else {
      await store.create({
        name: formName.value.trim(),
        parent_id: parentCategory.value?.id || null,
      });
      ElMessage.success("创建成功");
    }
    dialogVisible.value = false;
  } catch (err: any) {
    ElMessage.error(err.message || "操作失败");
  } finally {
    saving.value = false;
  }
}

async function handleDelete(data: Category) {
  try {
    await ElMessageBox.confirm(
      `确定要删除分类「${data.name}」吗？子分类也会被删除。`,
      "确认删除",
      { type: "warning" }
    );
    await store.remove(data.id);
    ElMessage.success("删除成功");
  } catch {
    // cancelled
  }
}

function allowDrop(draggingNode: any, dropNode: any, type: string) {
  // 不能拖到自己或子节点上
  if (draggingNode.data.id === dropNode.data.id) return false;
  if (type === "inner") {
    // 检查是否是自己的子节点
    let parent = dropNode.parent;
    while (parent) {
      if (parent.data.id === draggingNode.data.id) return false;
      parent = parent.parent;
    }
  }
  return true;
}

async function handleDrop(draggingNode: any, dropNode: any, dropType: string) {
  const dragId = draggingNode.data.id;
  let parentId: number | null = null;

  if (dropType === "inner") {
    parentId = dropNode.data.id;
  } else {
    // before/after: same parent as dropNode
    parentId = dropNode.data.parent_id;
  }

  try {
    await store.move(dragId, {
      parent_id: parentId,
      sort_order: 0,
    });
    ElMessage.success("移动成功");
  } catch (err: any) {
    ElMessage.error("移动失败");
    store.fetchTree(); // 刷新以恢复
  }
}

// 初始化加载
store.fetchTree();
</script>

<style scoped>
.tree-node-content {
  display: flex;
  align-items: center;
  flex: 1;
  overflow: hidden;
}
.tree-node-label {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
