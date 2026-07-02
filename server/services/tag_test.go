package services

import (
	"testing"

	"pcms/store"
)

func TestTagService_List(t *testing.T) {
	s := NewTagService(store.NewMemStore())

	// 空列表
	tags, err := s.List(1)
	if err != nil {
		t.Fatalf("期望成功, 得到: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("空列表应有 0 个, 得到 %d", len(tags))
	}
}

func TestTagService_Create(t *testing.T) {
	s := NewTagService(store.NewMemStore())

	t.Run("创建标签", func(t *testing.T) {
		tag, err := s.Create(1, CreateTagInput{Name: "Go"})
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if tag.Name != "Go" {
			t.Fatalf("期望 Go, 得到 %s", tag.Name)
		}
	})

	t.Run("去重", func(t *testing.T) {
		tag, err := s.Create(1, CreateTagInput{Name: "Go"})
		if err != nil {
			t.Fatalf("FirstOrCreate 应成功, 得到: %v", err)
		}
		if tag.ID != 1 {
			t.Fatalf("应返回已有标签 ID=1, 得到 %d", tag.ID)
		}
	})

	t.Run("用户隔离", func(t *testing.T) {
		tag, err := s.Create(2, CreateTagInput{Name: "Go"})
		if err != nil {
			t.Fatalf("不同用户应可创建同名标签, 得到: %v", err)
		}
		if tag.ID == 0 {
			t.Fatal("ID 不应为 0")
		}
	})
}

func TestTagService_Delete(t *testing.T) {
	s := NewTagService(store.NewMemStore())
	s.Create(1, CreateTagInput{Name: "Del"})

	t.Run("删除", func(t *testing.T) {
		err := s.Delete(1, 1)
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}

		tags, _ := s.List(1)
		if len(tags) != 0 {
			t.Fatal("删除后应无标签")
		}
	})

	t.Run("标签不存在", func(t *testing.T) {
		err := s.Delete(1, 999)
		if err == nil {
			t.Fatal("期望报错")
		}
	})
}
