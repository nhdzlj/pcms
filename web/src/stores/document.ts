import { defineStore } from "pinia";
import { ref } from "vue";
import {
  listDocuments,
  searchDocuments,
  getDocument,
  createDocument,
  updateDocument,
  deleteDocument,
} from "@/api/document";
import type { Document, CreateDocumentParams, UpdateDocumentParams, PaginatedResult } from "@/api/document";

export const useDocumentStore = defineStore("document", () => {
  const documents = ref<Document[]>([]);
  const current = ref<Document | null>(null);
  const total = ref(0);
  const loading = ref(false);
  const currentPage = ref(1);
  const pageSize = ref(20);

  async function fetchList(params?: {
    page?: number;
    category_id?: number;
    status?: string;
  }) {
    loading.value = true;
    try {
      const page = params?.page || 1;
      const result: PaginatedResult<Document> = await listDocuments({
        page,
        page_size: pageSize.value,
        category_id: params?.category_id,
        status: params?.status,
      });
      documents.value = result.list;
      total.value = result.pagination.total;
      currentPage.value = result.pagination.page;
    } finally {
      loading.value = false;
    }
  }

  async function search(keyword: string, page = 1) {
    loading.value = true;
    try {
      const result = await searchDocuments(keyword, page, pageSize.value);
      documents.value = result.list;
      total.value = result.pagination.total;
      currentPage.value = result.pagination.page;
    } finally {
      loading.value = false;
    }
  }

  async function fetchById(id: number) {
    loading.value = true;
    try {
      current.value = await getDocument(id);
    } finally {
      loading.value = false;
    }
  }

  async function create(params: CreateDocumentParams) {
    const doc = await createDocument(params);
    return doc;
  }

  async function update(id: number, params: UpdateDocumentParams) {
    const doc = await updateDocument(id, params);
    current.value = doc;
    return doc;
  }

  async function remove(id: number) {
    await deleteDocument(id);
    current.value = null;
    await fetchList();
  }

  return {
    documents,
    current,
    total,
    loading,
    currentPage,
    pageSize,
    fetchList,
    search,
    fetchById,
    create,
    update,
    remove,
  };
});
