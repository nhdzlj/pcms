package services

import (
	"testing"

	"pcms/store"
	"pcms/utils"
)

func TestDocumentService_Create(t *testing.T) {
	utils.InitJWT("test-secret")
	s := NewDocumentService(store.NewMemStore())

	doc, err := s.Create(1, CreateDocumentInput{
		Title:   "Test Doc",
		Content: "Hello",
	})
	if err != nil {
		t.Fatalf("期望成功, 得到: %v", err)
	}
	if doc.ID != 1 {
		t.Fatalf("期望 ID=1, 得到 %d", doc.ID)
	}
	if doc.Version != 1 {
		t.Fatalf("期望 Version=1, 得到 %d", doc.Version)
	}
	if doc.Status != "draft" {
		t.Fatalf("默认 status=draft, 得到 %s", doc.Status)
	}
}

func TestDocumentService_List(t *testing.T) {
	s := NewDocumentService(store.NewMemStore())
	s.Create(1, CreateDocumentInput{Title: "A", Content: "a"})
	s.Create(1, CreateDocumentInput{Title: "B", Content: "b"})

	result, err := s.List(1, 1, 20, nil, "", nil)
	if err != nil {
		t.Fatalf("期望成功, 得到: %v", err)
	}
	if result.Pagination.Total != 2 {
		t.Fatalf("期望 2, 得到 %d", result.Pagination.Total)
	}
}

func TestDocumentService_GetByID(t *testing.T) {
	s := NewDocumentService(store.NewMemStore())
	s.Create(1, CreateDocumentInput{Title: "My Doc", Content: "c"})

	doc, err := s.GetByID(1, 1)
	if err != nil {
		t.Fatalf("期望成功, 得到: %v", err)
	}
	if doc.Title != "My Doc" {
		t.Fatalf("期望 My Doc, 得到 %s", doc.Title)
	}
}

func TestDocumentService_Update(t *testing.T) {
	s := NewDocumentService(store.NewMemStore())
	s.Create(1, CreateDocumentInput{Title: "Old", Content: "old"})

	doc, err := s.Update(1, 1, UpdateDocumentInput{
		Title:   "New",
		Content: "new",
	})
	if err != nil {
		t.Fatalf("期望成功, 得到: %v", err)
	}
	if doc.Title != "New" {
		t.Fatalf("期望 New, 得到 %s", doc.Title)
	}
	if doc.Version != 2 {
		t.Fatalf("期望 Version=2, 得到 %d", doc.Version)
	}
}

func TestDocumentService_Delete(t *testing.T) {
	s := NewDocumentService(store.NewMemStore())
	s.Create(1, CreateDocumentInput{Title: "Del", Content: "d"})

	err := s.Delete(1, 1)
	if err != nil {
		t.Fatalf("期望成功, 得到: %v", err)
	}

	err = s.Delete(1, 1)
	if err == nil {
		t.Fatal("期望报错")
	}
}

func TestDocumentService_UserIsolation(t *testing.T) {
	s := NewDocumentService(store.NewMemStore())
	s.Create(1, CreateDocumentInput{Title: "Mine", Content: "c"})

	_, err := s.GetByID(2, 1)
	if err == nil {
		t.Fatal("用户 2 不应访问用户 1 的文档")
	}

	_, err = s.Update(2, 1, UpdateDocumentInput{Title: "x", Content: "y"})
	if err == nil {
		t.Fatal("用户 2 不应更新用户 1 的文档")
	}
}
