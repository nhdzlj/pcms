import { defineStore } from "pinia";
import { ref } from "vue";
import {
  listDocuments,
  searchDocuments,
  getDocument,
  createDocument,
  updateDocument,
  deleteDocument,
  getDocumentVersions,
  getDocumentVersion,
} from "@/api/document";
import type { Document, CreateDocumentParams, UpdateDocumentParams, PaginatedResult, DocumentVersion } from "@/api/document";

export const useDocumentStore = defineStore("document", () => {
  const documents = ref<Document[]>([]);
  const current = ref<Document | null>(null);
  const total = ref(0);
  const loading = ref(false);
  const currentPage = ref(1);
  const pageSize = ref(20);
  const versions = ref<DocumentVersion[]>([]);
  const currentVersion = ref<DocumentVersion | null>(null);

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

  async function search(keyword: string, page = 1, params?: { tag_id?: number; category_id?: number }) {
    loading.value = true;
    try {
      const result = await searchDocuments({
        keyword,
        page,
        page_size: pageSize.value,
        ...params,
      });
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

  async function fetchVersions(documentId: number) {
    versions.value = await getDocumentVersions(documentId);
  }

  async function fetchVersion(documentId: number, versionId: number) {
    currentVersion.value = await getDocumentVersion(documentId, versionId);
    return currentVersion.value;
  }

  return {
    documents,
    current,
    total,
    loading,
    currentPage,
    pageSize,
    versions,
    currentVersion,
    fetchList,
    search,
    fetchById,
    create,
    update,
    remove,
    fetchVersions,
    fetchVersion,
  };
});
