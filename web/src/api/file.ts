import request from "./request";

export interface UploadResult {
  url: string;
  file_name: string;
  file_size: number;
  mime_type: string;
}

export function uploadFile(file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return request.post<any, UploadResult>("/files/upload", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
}
