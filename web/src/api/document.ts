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
}

export function listDocuments(params?: ListDocumentsParams) {
  return request.get<any, PaginatedResult<Document>>("/documents", { params });
}

export function searchDocuments(keyword: string, page = 1, pageSize = 20) {
  return request.get<any, PaginatedResult<Document>>("/documents/search", {
    params: { keyword, page, page_size: pageSize },
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
