import request from "./request";

export interface Attachment {
  id: number;
  document_id: number | null;
  file_name: string;
  file_path: string;
  file_size: number;
  mime_type: string;
  created_at: string;
}

export function getAttachments(params?: { document_id?: number }) {
  return request.get<any, Attachment[]>("/attachments", { params });
}

export function createAttachment(
  params: {
    file_name: string;
    file_path: string;
    file_size: number;
    mime_type: string;
  },
  documentId?: number
) {
  return request.post<any, Attachment>(
    "/attachments",
    params,
    { params: documentId ? { document_id: documentId } : undefined }
  );
}

export function deleteAttachment(id: number) {
  return request.delete(`/attachments/${id}`);
}

export function bindAttachmentToDocument(id: number, documentId: number) {
  return request.put(`/attachments/${id}/bind`, { document_id: documentId });
}
