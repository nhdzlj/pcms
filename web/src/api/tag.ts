import request from "./request";

export interface Tag {
  id: number;
  name: string;
}

export function getTags() {
  return request.get<any, Tag[]>("/tags");
}

export function createTag(name: string) {
  return request.post<any, Tag>("/tags", { name });
}

export function deleteTag(id: number) {
  return request.delete(`/tags/${id}`);
}
