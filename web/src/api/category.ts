import request from "./request";

export interface Category {
  id: number;
  name: string;
  parent_id: number | null;
  sort_order: number;
  icon: string;
  children: Category[];
  created_at: string;
  updated_at: string;
}

export interface CreateCategoryParams {
  name: string;
  parent_id?: number | null;
  icon?: string;
}

export interface UpdateCategoryParams {
  name: string;
  parent_id?: number | null;
  sort_order?: number;
  icon?: string;
}

export interface MoveCategoryParams {
  parent_id: number | null;
  sort_order: number;
}

export function getCategoryTree() {
  return request.get<any, Category[]>("/categories");
}

export function createCategory(params: CreateCategoryParams) {
  return request.post<any, Category>("/categories", params);
}

export function updateCategory(id: number, params: UpdateCategoryParams) {
  return request.put<any, Category>(`/categories/${id}`, params);
}

export function moveCategory(id: number, params: MoveCategoryParams) {
  return request.put<any, Category>(`/categories/${id}/move`, params);
}

export function deleteCategory(id: number) {
  return request.delete(`/categories/${id}`);
}
