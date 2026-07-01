import { defineStore } from "pinia";
import { ref } from "vue";
import {
  getCategoryTree,
  createCategory,
  updateCategory,
  deleteCategory,
  moveCategory,
} from "@/api/category";
import type { Category, CreateCategoryParams, UpdateCategoryParams, MoveCategoryParams } from "@/api/category";

export const useCategoryStore = defineStore("category", () => {
  const tree = ref<Category[]>([]);
  const loading = ref(false);
  const selectedId = ref<number | null>(null);

  async function fetchTree() {
    loading.value = true;
    try {
      tree.value = await getCategoryTree();
    } finally {
      loading.value = false;
    }
  }

  async function create(params: CreateCategoryParams) {
    await createCategory(params);
    await fetchTree();
  }

  async function update(id: number, params: UpdateCategoryParams) {
    await updateCategory(id, params);
    await fetchTree();
  }

  async function remove(id: number) {
    await deleteCategory(id);
    if (selectedId.value === id) {
      selectedId.value = null;
    }
    await fetchTree();
  }

  async function move(id: number, params: MoveCategoryParams) {
    await moveCategory(id, params);
    await fetchTree();
  }

  function select(id: number | null) {
    selectedId.value = id;
  }

  return {
    tree,
    loading,
    selectedId,
    fetchTree,
    create,
    update,
    remove,
    move,
    select,
  };
});
