import request from "./request";

export interface Document {
  id: number;
  title: string;
  content: string;
  summary: string;
  category_id: number | null;
  status: string;
  view_count: number;
  is_favorite: boolean;
  version: number;
  category?: {
    id: number;
    name: string;
  };
  tags: { id: number; name: string }[];
  created_at: string;
  updated_at: string;
}

export interface PaginatedResult<T> {
  list: T[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
  };
}

export interface CreateDocumentParams {
  title: string;
  content: string;
  summary?: string;
  category_id?: number | null;
  tag_ids?: number[];
  status?: string;
}

export interface UpdateDocumentParams extends CreateDocumentParams {
  is_favorite?: boolean;
}

export interface ListDocumentsParams {
  page?: number;
  page_size?: number;
  category_id?: number;
  status?: string;
  tag_id?: number;
}

export interface SearchDocumentsParams {
  keyword: string;
  page?: number;
  page_size?: number;
  tag_id?: number;
  category_id?: number;
}

export function listDocuments(params?: ListDocumentsParams) {
  return request.get<any, PaginatedResult<Document>>("/documents", { params });
}

export function searchDocuments(params: SearchDocumentsParams) {
  return request.get<any, PaginatedResult<Document>>("/documents/search", {
    params,
  });
}

export function getDocument(id: number) {
  return request.get<any, Document>(`/documents/${id}`);
}

export function createDocument(params: CreateDocumentParams) {
  return request.post<any, Document>("/documents", params);
}

export function updateDocument(id: number, params: UpdateDocumentParams) {
  return request.put<any, Document>(`/documents/${id}`, params);
}

export function deleteDocument(id: number) {
  return request.delete(`/documents/${id}`);
}

// 文档版本相关
export interface DocumentVersion {
  id: number;
  document_id: number;
  version: number;
  title: string;
  content: string;
  created_at: string;
}

export function getDocumentVersions(documentId: number) {
  return request.get<any, DocumentVersion[]>(`/documents/${documentId}/versions`);
}

export function getDocumentVersion(documentId: number, versionId: number) {
  return request.get<any, DocumentVersion>(`/documents/${documentId}/versions/${versionId}`);
}
