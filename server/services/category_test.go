package services

import (
	"testing"

	"pcms/store"
)

func TestCategoryService_GetTree(t *testing.T) {
	s := NewCategoryService(store.NewMemStore())

	// 空树
	tree, err := s.GetTree(1)
	if err != nil {
		t.Fatalf("期望成功, 得到: %v", err)
	}
	if len(tree) != 0 {
		t.Fatalf("空树应有 0 个节点, 得到 %d", len(tree))
	}
}

func TestCategoryService_Create(t *testing.T) {
	s := NewCategoryService(store.NewMemStore())

	t.Run("创建根分类", func(t *testing.T) {
		cat, err := s.Create(1, CreateCategoryInput{Name: "Root"})
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if cat.Name != "Root" {
			t.Fatalf("期望 Root, 得到 %s", cat.Name)
		}
		if cat.Icon != "folder" {
			t.Fatalf("默认图标应为 folder, 得到 %s", cat.Icon)
		}
	})

	t.Run("创建子分类", func(t *testing.T) {
		parentID := uint64(1)
		cat, err := s.Create(1, CreateCategoryInput{
			Name:     "Child",
			ParentID: &parentID,
		})
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if cat.ParentID == nil || *cat.ParentID != 1 {
			t.Fatal("ParentID 应为 1")
		}
	})

	t.Run("父分类不存在", func(t *testing.T) {
		parentID := uint64(999)
		_, err := s.Create(1, CreateCategoryInput{
			Name:     "Orphan",
			ParentID: &parentID,
		})
		if err == nil {
			t.Fatal("期望报错")
		}
	})
}

func TestCategoryService_Update(t *testing.T) {
	s := NewCategoryService(store.NewMemStore())
	s.Create(1, CreateCategoryInput{Name: "Old"})

	t.Run("更新名称", func(t *testing.T) {
		cat, err := s.Update(1, 1, UpdateCategoryInput{Name: "New"})
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if cat.Name != "New" {
			t.Fatalf("期望 New, 得到 %s", cat.Name)
		}
	})

	t.Run("不能移到自身", func(t *testing.T) {
		selfID := uint64(1)
		_, err := s.Update(1, 1, UpdateCategoryInput{
			Name:     "Self",
			ParentID: &selfID,
		})
		if err == nil {
			t.Fatal("期望报错")
		}
	})

	t.Run("分类不存在", func(t *testing.T) {
		_, err := s.Update(1, 999, UpdateCategoryInput{Name: "X"})
		if err == nil {
			t.Fatal("期望报错")
		}
	})
}

func TestCategoryService_Delete(t *testing.T) {
	s := NewCategoryService(store.NewMemStore())
	s.Create(1, CreateCategoryInput{Name: "Del"})

	t.Run("删除", func(t *testing.T) {
		err := s.Delete(1, 1)
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
	})

	t.Run("不存在", func(t *testing.T) {
		err := s.Delete(1, 1)
		if err == nil {
			t.Fatal("期望报错")
		}
	})
}

func TestCategoryService_Move(t *testing.T) {
	s := NewCategoryService(store.NewMemStore())
	s.Create(1, CreateCategoryInput{Name: "A"})
	s.Create(1, CreateCategoryInput{Name: "B"})

	t.Run("移动", func(t *testing.T) {
		parentID := uint64(1)
		cat, err := s.Move(1, 2, MoveCategoryInput{
			ParentID:  &parentID,
			SortOrder: 1,
		})
		if err != nil {
			t.Fatalf("期望成功, 得到: %v", err)
		}
		if cat.ParentID == nil || *cat.ParentID != 1 {
			t.Fatal("ParentID 应为 1")
		}
	})

	t.Run("不能移到自身", func(t *testing.T) {
		selfID := uint64(2)
		_, err := s.Move(1, 2, MoveCategoryInput{
			ParentID:  &selfID,
			SortOrder: 0,
		})
		if err == nil {
			t.Fatal("期望报错")
		}
	})
}
