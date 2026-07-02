package services

import (
	"testing"

	"pcms/store"
)

func TestAttachmentService_Create(t *testing.T) {
	s := NewAttachmentService(store.NewMemStore())

	t.Run("创建", func(t *testing.T) {
		a, err := s.Create(1, CreateAttachmentInput{
			FileName: "test.pdf",
			FilePath: "/uploads/test.pdf",
			FileSize: 1024,
			MimeType: "application/pdf",
		}, nil)
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if a.ID != 1 {
			t.Fatalf("期望 ID=1, 得到 %d", a.ID)
		}
	})

	t.Run("关联文档", func(t *testing.T) {
		docID := uint64(10)
		a, err := s.Create(1, CreateAttachmentInput{
			FileName: "linked.pdf",
			FilePath: "/uploads/linked.pdf",
		}, &docID)
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if a.DocumentID == nil || *a.DocumentID != 10 {
			t.Fatal("DocumentID 应为 10")
		}
	})
}

func TestAttachmentService_List(t *testing.T) {
	s := NewAttachmentService(store.NewMemStore())
	s.Create(1, CreateAttachmentInput{FileName: "a.png", FilePath: "/a.png"}, nil)
	s.Create(1, CreateAttachmentInput{FileName: "b.png", FilePath: "/b.png"}, nil)

	list, err := s.List(1, ListAttachmentQuery{})
	if err != nil {
		t.Fatalf("期望成功, 得到: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("期望 2 个, 得到 %d", len(list))
	}
}

func TestAttachmentService_Delete(t *testing.T) {
	s := NewAttachmentService(store.NewMemStore())
	s.Create(1, CreateAttachmentInput{FileName: "d.pdf", FilePath: "/d.pdf"}, nil)

	err := s.Delete(1, 1)
	if err != nil {
		t.Fatalf("期望成功, 得到: %v", err)
	}

	err = s.Delete(1, 1)
	if err == nil {
		t.Fatal("期望报错")
	}
}
